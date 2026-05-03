import argparse
import os
from dotenv import load_dotenv
from langchain_community.document_loaders import PyPDFLoader
from langchain_text_splitters import CharacterTextSplitter
from langchain_openai import OpenAIEmbeddings
from langchain_qdrant import QdrantVectorStore

load_dotenv()

parser = argparse.ArgumentParser()
parser.add_argument("path", help="Path to the PDF file to index")
args = parser.parse_args()

loader = PyPDFLoader(args.path)
docs = loader.load()

splitter = CharacterTextSplitter(chunk_size=1000, chunk_overlap=100)
chunks = splitter.split_documents(docs)

embeddings = OpenAIEmbeddings(model="text-embedding-3-small")

QdrantVectorStore.from_documents(
    chunks,
    embeddings,
    url=os.getenv("QDRANT_URL"),
    collection_name="documents",
)

print(f"Indexed {len(chunks)} chunks from {os.path.basename(args.path)}")
