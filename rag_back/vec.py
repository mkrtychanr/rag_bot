import os
import asyncio
import json
import re
from pathlib import Path
from fastapi import FastAPI, HTTPException
from pymilvus import connections, FieldSchema, CollectionSchema, DataType, Collection, utility
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.schema import Document
from langchain_community.document_loaders import TextLoader
from openai import OpenAI, AsyncOpenAI
import openai
from env import OPENAI_API

from langchain_openai import OpenAIEmbeddings
from tenacity import retry, wait_random_exponential, stop_after_attempt

# Конфигурация
MILVUS_HOST = '0.0.0.0'
MILVUS_PORT = '19530'
COLLECTION_NAME = 'rag_documents_lab_gffg'
EMBEDDING_MODEL = "text-embedding-3-large"
DIMENSION = 3072

# Настройка OpenAI
llm = AsyncOpenAI(api_key=OPENAI_API)
openai.api_key = OPENAI_API

class Document:
    def __init__(self, page_content, doc_id=None):
        self.page_content = page_content
        self.doc_id = doc_id  # Теперь metadata заменено на doc_id

class DocumentLoader:
    def __init__(self, file_path=None, path=None):
        self.file_path = file_path
        self.path = path

    def load_documents(self, doc_id=None):
        """Загружает документы, принимая doc_id вручную"""
        documents = []
        if self.file_path:
            try:
                with open(self.file_path, 'r', encoding='utf-8') as file:
                    content = file.read()
                    documents.append(Document(
                        page_content=content,
                        doc_id=doc_id  # Используем переданный doc_id
                    ))
            except Exception as e:
                print(f"Ошибка при загрузке файла {self.file_path}: {str(e)}")
        elif self.path:
            for filename in os.listdir(self.path):
                if filename.lower().endswith('.txt'):
                    file_path = os.path.join(self.path, filename)
                    try:
                        with open(file_path, 'r', encoding='utf-8') as file:
                            content = file.read()
                            # Для каждого файла нужно запросить doc_id
                            documents.append(Document(
                                page_content=content,
                                doc_id=doc_id  # Используем переданный doc_id
                            ))
                    except Exception as e:
                        print(f"Ошибка при загрузке файла {file_path}: {str(e)}")
        return documents

class Vectorestore:
    @retry(wait=wait_random_exponential(min=1, max=20), stop=stop_after_attempt(3))
    def get_embeddings(self, texts):
        response = openai.embeddings.create(
            model=EMBEDDING_MODEL,
            input=texts
        )
        return [data.embedding for data in response.data]

    async def connect_to_milvus(self):
        try:
            connections.connect(
                alias='default',
                host=MILVUS_HOST,
                port=MILVUS_PORT,
            )
        except Exception as e:
            print(f"Ошибка подключения к Milvus: {str(e)}")
            raise
    async def ensure_collection_exists(self):
        """Проверяет существование коллекции и создаёт её при необходимости"""
        await self.connect_to_milvus()
        if not utility.has_collection(COLLECTION_NAME):
            self.create_collection()
    async def delete_document(self, doc_id: int):
        """Удаляет все чанки документа по его ID"""
        await self.ensure_collection_exists()
        collection = Collection(COLLECTION_NAME)
        collection.load()
        
        expr = f"metadata == {doc_id}"
        try:
            res = collection.delete(expr=expr)
            return res
        except Exception as e:
            print(f"Ошибка при удалении документа {doc_id}: {str(e)}")
            raise
    def create_collection(self):
        if utility.has_collection(COLLECTION_NAME):
            utility.drop_collection(COLLECTION_NAME)
        fields = [
            FieldSchema(name='id', dtype=DataType.INT64, is_primary=True, auto_id=True),
            FieldSchema(name='text', dtype=DataType.VARCHAR, max_length=2000),
            FieldSchema(name='metadata', dtype=DataType.INT64),  # Теперь храним только doc_id
            FieldSchema(name='embedding', dtype=DataType.FLOAT_VECTOR, dim=DIMENSION),
        ]
        schema = CollectionSchema(fields=fields, description='RAG documents storage')
        collection = Collection(name=COLLECTION_NAME, schema=schema)
        index_params = {
            'metric_type': 'IP',
            'index_type': 'IVF_FLAT',
            'params': {'nlist': 256}
        }
        collection.create_index(field_name='embedding', index_params=index_params)
        return collection

    def index_documents(self, documents, batch_size=10):
        asyncio.run(self.connect_to_milvus())
        if not utility.has_collection(COLLECTION_NAME):
            self.create_collection()
        collection = Collection(COLLECTION_NAME)

        def clean_text(text):
            text = re.sub(r"<.*?>", "", text)
            text = re.sub(r"[^\w\s]", "", text)
            return text.strip()

        def preprocess_chunks(chunks):
            return [clean_text(chunk) for chunk in chunks if clean_text(chunk)]

        text_splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=100)
        all_chunks = []
        all_doc_ids = []

        for document in documents:
            if hasattr(document, "page_content"):
                text = document.page_content
                if isinstance(text, str):
                    chunks = text_splitter.split_text(text)
                    doc_id = document.doc_id  # Получаем doc_id из документа
                    if doc_id is None:
                        raise ValueError("doc_id должен быть указан для каждого документа")
                    all_chunks.extend(chunks)
                    all_doc_ids.extend([doc_id] * len(chunks))  # Используем doc_id вместо metadata
                else:
                    raise ValueError(f"Содержимое документа должно быть строкой, но получен тип {type(text)}")
            else:
                raise ValueError("Объект документа не содержит атрибута page_content")

        for i in range(0, len(all_chunks), batch_size):
            batch_chunks = all_chunks[i:i + batch_size]
            batch_doc_ids = all_doc_ids[i:i + batch_size]
            batch_chunks = preprocess_chunks(batch_chunks)
            try:
                embeddings = self.get_embeddings(batch_chunks)
                data = [batch_chunks, batch_doc_ids, embeddings]  # Теперь передаем doc_ids вместо metadata
                collection.insert(data)
            except Exception as e:
                print(f"Ошибка при вставке батча {i//batch_size + 1}: {str(e)}")

        collection.load()
        print(f"Успешно проиндексировано {len(all_chunks)} чанков")

    async def search(self, query, top_k=10, document_ids=None):
        """Поиск с возможностью фильтрации по списку document_ids"""
        await self.ensure_collection_exists()
        collection = Collection(COLLECTION_NAME)
        print(collection.load())

        query_embedding = self.get_embeddings([query])[0]

        search_params = {
            'metric_type': 'IP',
            'params': {'nprobe': 32}
        }

        # Создаем выражение для фильтрации по document_ids
        expr = None
        if document_ids:
            if len(document_ids) == 1:
                expr = f"metadata == {document_ids[0]}"
            else:
                doc_ids_str = ", ".join(str(doc_id) for doc_id in document_ids)
                expr = f"metadata in [{doc_ids_str}]"

        try:
            results = collection.search(
                data=[query_embedding],
                anns_field='embedding',
                param=search_params,
                limit=top_k,
                output_fields=['text', 'metadata'],
                expr=expr  # Применяем фильтр по document_ids
            )
        except Exception as e:
            print(f"Ошибка при поиске: {str(e)}")
            raise

        context = []
        for hit in results[0]:
            text = hit.entity.get('text')
            doc_id = hit.entity.get('metadata')
            context.append({
                "text": text,
                "doc_id": doc_id
            })
        return context

    async def check_existing_documents(self, texts, similarity_threshold=0.9):
        await self.connect_to_milvus()
        collection = Collection(COLLECTION_NAME)
        collection.load()

        embeddings_to_check = self.get_embeddings(texts)

        search_params = {
            'metric_type': 'IP',
            'params': {'nprobe': 32}
        }

        try:
            results = collection.search(
                data=embeddings_to_check,
                anns_field='embedding',
                param=search_params,
                limit=1,
                expr=None
            )
        except Exception as e:
            print(f"Ошибка при проверке существующих документов: {str(e)}")
            return []

        existing_texts = []
        for i, hits in enumerate(results):
            if hits and hits[0].distance > similarity_threshold:
                existing_texts.append(texts[i])

        return existing_texts

    async def append_documents(self, new_documents, batch_size=10):
        await self.connect_to_milvus()
        collection = Collection(COLLECTION_NAME)

        text_splitter = RecursiveCharacterTextSplitter(chunk_size=900, chunk_overlap=100)
        all_chunks = []
        all_doc_ids = []
        texts_to_check = []

        for document in new_documents:
            if hasattr(document, "page_content"):
                text = document.page_content
                if isinstance(text, str):
                    chunks = text_splitter.split_text(text)
                    doc_id = document.doc_id  # Получаем doc_id из документа
                    if doc_id is None:
                        raise ValueError("doc_id должен быть указан для каждого документа")
                    all_chunks.extend(chunks)
                    all_doc_ids.extend([doc_id] * len(chunks))  # Используем doc_id вместо metadata
                    texts_to_check.extend(chunks)
                else:
                    raise ValueError(f"Содержимое документа должно быть строкой, но получен тип {type(text)}")
            else:
                raise ValueError("Объект документа не содержит атрибута page_content")

        try:
            existing_texts = await self.check_existing_documents(texts_to_check)
        except Exception as e:
            print(f"Ошибка при проверке дубликатов: {str(e)}")
            existing_texts = []

        new_chunks = [chunk for chunk in all_chunks if chunk not in existing_texts]
        new_doc_ids = [doc_id for chunk, doc_id in zip(all_chunks, all_doc_ids) if chunk not in existing_texts]

        for i in range(0, len(new_chunks), batch_size):
            batch_chunks = new_chunks[i:i + batch_size]
            batch_doc_ids = new_doc_ids[i:i + batch_size]
            try:
                embeddings = self.get_embeddings(batch_chunks)
                data = [batch_chunks, batch_doc_ids, embeddings]  # Теперь передаем doc_ids вместо metadata
                collection.insert(data)
            except Exception as e:
                print(f"Ошибка при вставке батча {i//batch_size + 1}: {str(e)}")
                continue

        collection.load()
        print(f"Успешно добавлено {len(new_chunks)} новых чанков")

    async def generate_answer(self, query, context, document_ids=None, model="gpt-4o-mini"):
        """Генерация ответа с возможностью фильтрации по document_ids"""
        # Фильтруем контекст по document_ids, если они указаны
        if document_ids is not None:
            context = [item for item in context if item['doc_id'] in document_ids]
            if not context:
                return "Нет информации по запросу в указанных документах"

        system_prompt = """
Ты — эксперт в области корпоративной базы знаний. Твоя задача — предоставить информацию, которая будет полезна для сотрудников. Используй только данные из предоставленных документов, избегай домыслов и неподтвержденной информации.
"""
        texts = [item['text'] for item in context]
        user_content = f"Контекст:\n{'\n'.join(texts)}\n\nТема: {query}"

        try:
            response = await llm.chat.completions.create(
                model=model,
                messages=[
                    {"role": "system", "content": system_prompt},
                    {"role": "user", "content": user_content}
                ],
                temperature=0.1,
            )
            return response.choices[0].message.content
        except Exception as e:
            raise HTTPException(status_code=500, detail=f"Error generating text")

    async def get_all_doc_ids(self):
        """Возвращает все уникальные doc_id из коллекции"""
        await self.connect_to_milvus()
        collection = Collection(COLLECTION_NAME)
        collection.load()

        expr = "metadata != 0"  # Предполагаем, что doc_id не может быть 0
        output_fields = ['metadata']

        try:
            results = collection.query(
                expr=expr,
                output_fields=output_fields
            )
        except Exception as e:
            print(f"Ошибка при получении doc_id: {str(e)}")
            return []

        doc_ids = set()
        for result in results:
            doc_ids.add(result['metadata'])

        return list(doc_ids)

if __name__ == "__main__":
    async def main():
        #### Пример загрузки нового документа с указанием doc_id
        vector_store = Vectorestore()
        await vector_store.connect_to_milvus()
        
        # Создаем коллекцию (если нужно)
        # vector_store.create_collection()
        
        # Загружаем несколько документов с разными doc_id
        doc1_id = 12345
        doc2_id = 67890
        
        # # Загрузка первого документа
        # document_loader1 = DocumentLoader(file_path="/Users/thetom205/raffi_rag/RAF_DOC/Лабораторная работа_4.txt")
        # new_documents1 = document_loader1.load_documents(doc_id=doc1_id)
        # await vector_store.append_documents(new_documents1)
        
        # # Загрузка второго документа
        # document_loader2 = DocumentLoader(file_path="/Users/thetom205/raffi_rag/RAF_DOC/НЕЙРОВОСПАЛЕНИЕ 1 (1)_008.txt")
        # new_documents2 = document_loader2.load_documents(doc_id=doc2_id)
        # await vector_store.append_documents(new_documents2)
        
        # Получаем все doc_id в коллекции
        doc_ids = await vector_store.get_all_doc_ids()
        print(f"Текущие doc_id в коллекции: {doc_ids}")
        
        #### Тестовый запрос ко всем документам
        query = "Расскажи про Ю 116"

        #### Тестовый запрос только к определенному документу
        print("\nЗапрос только к документу с ID", doc2_id)
        try:
            # Ищем только в документе с doc1_id
            context = await vector_store.search(query, top_k=10, document_ids=[doc2_id,doc1_id])
            answer = await vector_store.generate_answer(query, context, document_ids=[doc2_id,doc1_id])
            print("Ответ:", answer)
        except Exception as e:
            print(f"Ошибка при поиске или генерации ответа: {str(e)}")

    asyncio.run(main())