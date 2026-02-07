import sys
from pathlib import Path
from fastapi.middleware.cors import CORSMiddleware
from core.config import get_settings

sys.path.append(str(Path(__file__).resolve().parent))

import uvicorn
from fastapi import FastAPI
from core.lifespan import lifespan
from api.v1.routers import api_router as v1_router

settings = get_settings()


app = FastAPI(lifespan=lifespan, title=settings.PROJECT_NAME)
app.include_router(v1_router, prefix=settings.API_V1_PREFIX)


origins = settings.CORS_ORIGINS.split(",")
app.add_middleware(
  CORSMiddleware,
  allow_origins=origins,
  allow_credentials=True,
  allow_methods=["*"],
  allow_headers=["*"],
)

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
