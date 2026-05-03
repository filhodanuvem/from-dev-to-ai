from dotenv import load_dotenv
import argparse
from langchain_core.tools import tool
from langchain_openai import ChatOpenAI
import requests

load_dotenv()


@tool
def get_monster_type(id_or_name: str) -> list[str]:
    """Returns the types of a monster given its id or name."""
    response = requests.get(f"https://pokeapi.co/api/v2/pokemon/{id_or_name}/")
    response.raise_for_status()
    data = response.json()
    return [entry["type"]["name"] for entry in data["types"]]

parser = argparse.ArgumentParser()
parser.add_argument("query", help="Question to answer about")
args = parser.parse_args()

llm = ChatOpenAI(model="gpt-4o-mini", temperature=0)
llm_with_tools = llm.bind_tools([get_monster_type])

tools_by_name = {get_monster_type.name: get_monster_type}

response = llm_with_tools.invoke(args.query)

if response.tool_calls:
    for tool_call in response.tool_calls:
        tool = tools_by_name[tool_call["name"]]
        result = tool.invoke(tool_call["args"])
        print(result)
else:
    print(response.content)
