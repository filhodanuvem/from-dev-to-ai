from langchain_openai import ChatOpenAI
from langchain_google_genai import ChatGoogleGenerativeAI
from dotenv import load_dotenv


load_dotenv()

# llm = ChatOpenAI(
#   model="gpt-4o-mini",
#   temperature=0
# )

llm = ChatGoogleGenerativeAI(
  model="gemini-2.5-flash",
  temperature=0
)

response = llm.invoke("What is the capital of Portugal?")
print(response.content)