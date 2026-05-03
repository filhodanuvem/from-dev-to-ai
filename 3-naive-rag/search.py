import argparse
import os
from dotenv import load_dotenv
from langchain_openai import OpenAIEmbeddings, ChatOpenAI
from langchain_qdrant import QdrantVectorStore

load_dotenv()

parser = argparse.ArgumentParser()
parser.add_argument("query", help="Question to answer from the indexed documents")
args = parser.parse_args()

embeddings = OpenAIEmbeddings(model="text-embedding-3-small")

vector_store = QdrantVectorStore.from_existing_collection(
    embedding=embeddings,
    url=os.getenv("QDRANT_URL"),
    collection_name="documents",
)

docs = vector_store.similarity_search(args.query, k=4)
context = "\n\n".join(doc.page_content for doc in docs)

llm = ChatOpenAI(model="gpt-4o-mini", temperature=0)
prompt = f"Answer using only the context below.\n\nContext:\n{context}\n\nQuestion: {args.query}"
response = llm.invoke(prompt)

print(response.content)
