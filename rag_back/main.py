from fastapi import FastAPI, HTTPException, UploadFile, File, Query
from typing import List, Optional
import logging
import re
import os
from tempfile import NamedTemporaryFile
from pydantic import BaseModel
from vec import Vectorestore, DocumentLoader

app = FastAPI()
logger = logging.getLogger()
vecstore = Vectorestore()

class FactResponse(BaseModel):
    facts: str
    doc_ids: List[int]

class UploadResponse(BaseModel):
    status: str
    doc_id: int
    # chunks_count: int
from fastapi import FastAPI, HTTPException, UploadFile, File, Query
from typing import List, Optional
import logging
import re
import os
from tempfile import NamedTemporaryFile
from pydantic import BaseModel
import httpx
from urllib.parse import urlparse, parse_qs

# Ваши импорты
from vec import Vectorestore, DocumentLoader

app = FastAPI()
logger = logging.getLogger()
vecstore = Vectorestore()

class FactResponse(BaseModel):
    facts: str
    doc_ids: List[int]

class UploadResponse(BaseModel):
    status: str
    doc_id: int

class DriveLinkRequest(BaseModel):
    drive_link: str
    doc_id: Optional[int] = None
    overwrite: Optional[bool] = False

async def download_from_drive(drive_url: str, temp_file_path: str):
    """Скачивает файл по ссылке Google Drive и сохраняет его во временный файл."""
    try:
        # Извлекаем ID файла из URL
        parsed = urlparse(drive_url)
        if "drive.google.com" not in parsed.netloc:
            raise ValueError("Invalid Google Drive URL")

        file_id = None
        if "file/d/" in drive_url:
            file_id = drive_url.split("file/d/")[1].split("/")[0]
        elif "id=" in drive_url:
            file_id = parse_qs(parsed.query).get("id", [None])[0]

        if not file_id:
            raise ValueError("Could not extract file ID from Google Drive URL")

        download_url = f"https://drive.google.com/uc?export=download&id={file_id}"

        async with httpx.AsyncClient() as client:
            response = await client.get(download_url, follow_redirects=True)
            response.raise_for_status()

            with open(temp_file_path, "wb") as f:
                f.write(response.content)

        return True
    except Exception as e:
        logger.error(f"Failed to download from Google Drive: {str(e)}")
        if os.path.exists(temp_file_path):
            os.unlink(temp_file_path)
        raise

@app.post("/upload_from_drive/")
async def upload_from_drive(request: DriveLinkRequest):
    """Загружает файл по ссылке Google Drive."""
    try:
        await vecstore.ensure_collection_exists()
        existing_doc_ids = await vecstore.get_all_doc_ids()

        doc_id = request.doc_id if request.doc_id is not None else max(existing_doc_ids, default=0) + 1

        if doc_id in existing_doc_ids and not request.overwrite:
            raise HTTPException(
                status_code=400,
                detail=f"Document with id {doc_id} already exists. Set overwrite=True to replace."
            )

        with NamedTemporaryFile(delete=False) as temp_file:
            temp_path = temp_file.name

        await download_from_drive(request.drive_link, temp_path)

        document_loader = DocumentLoader(file_path=temp_path)
        documents = document_loader.load_documents(doc_id=doc_id)

        if doc_id in existing_doc_ids and request.overwrite:
            await vecstore.delete_document(doc_id)

        await vecstore.append_documents(documents)
        os.unlink(temp_path)

        return UploadResponse(
            status="success",
            doc_id=doc_id,
        )

    except Exception as e:
        logger.error(f"Error uploading document from Google Drive: {str(e)}")
        if "temp_path" in locals() and os.path.exists(temp_path):
            os.unlink(temp_path)
        raise HTTPException(status_code=500, detail=str(e))

# Остальные эндпоинты (/upload_document/, /gen_facts/, /health/) остаются без изменений
@app.post("/upload_document/")
async def upload_document(
    file: UploadFile = File(...),
    doc_id: int = None,
    overwrite: bool = False
):
    try:
        await vecstore.ensure_collection_exists()
        existing_doc_ids = await vecstore.get_all_doc_ids()
        if doc_id in existing_doc_ids and not overwrite:
            raise HTTPException(
                status_code=400,
                detail=f"Document with id {doc_id} already exists. Set overwrite=True to replace."
            )

        with NamedTemporaryFile(delete=False) as temp_file:
            content = await file.read()
            temp_file.write(content)
            temp_path = temp_file.name

        document_loader = DocumentLoader(file_path=temp_path)
        documents = document_loader.load_documents(doc_id=doc_id)
        
        if doc_id in existing_doc_ids and overwrite:
            await vecstore.delete_document(doc_id)
        
        await vecstore.append_documents(documents)
        
        # chunks_count = await vecstore.get_chunks_count(doc_id)
        os.unlink(temp_path)
        
        return UploadResponse(
            status="success",
            doc_id=doc_id,
            # chunks_count=chunks_count
        )
        
    except Exception as e:
        logger.error(f"Error uploading document: {str(e)}")
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/gen_facts/")
async def generate_facts(
    query: str,
    doc_ids: Optional[List[int]] = Query(None)
):
    try:
        print("before connect")
        await vecstore.connect_to_milvus()
        print("after connect")
        context = await vecstore.search(query, top_k=10, document_ids=doc_ids)
        answer = await vecstore.generate_answer(query, context, document_ids=doc_ids)
        
        def clean_text(text):
            cleaned_text = re.sub(r'[#*]', '', text)
            cleaned_text = re.sub(r'\s+', ' ', cleaned_text).strip()
            return cleaned_text
        
        context_doc_ids = list(set([item['doc_id'] for item in context])) if context else []
        
        return FactResponse(
            facts=clean_text(answer),
            doc_ids=context_doc_ids
        )
    except Exception as e:
        logger.error(e)
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/health")
async def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)