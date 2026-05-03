# Naive RAG

### Requirements

* [pyenv](https://github.com/pyenv/pyenv#installation)

### Problem

In this chapter the idea is to do a RAG but with a "local" vector database using postgres + pgvector or qdrant for example.
Also would be interesting to use langchain.

* create .env file from .env.template.
* pip install -r requirements.txt
* python indexing.py _path_to_pdf
* python search "any question"