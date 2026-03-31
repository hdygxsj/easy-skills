"""Configuration API."""
from typing import List
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from ..database import get_db
from ..schemas import ConfigUpdate, ConfigResponse
from ..models import Config

router = APIRouter()


@router.get("/", response_model=List[ConfigResponse])
def list_config(db: Session = Depends(get_db)):
    """List all configuration."""
    return db.query(Config).all()


@router.get("/{key}", response_model=ConfigResponse)
def get_config(key: str, db: Session = Depends(get_db)):
    """Get configuration value."""
    config = db.query(Config).filter(Config.key == key).first()
    if not config:
        return ConfigResponse(key=key, value=None, updated_at=None)
    return config


@router.put("/{key}", response_model=ConfigResponse)
def set_config(key: str, data: ConfigUpdate, db: Session = Depends(get_db)):
    """Set configuration value."""
    config = db.query(Config).filter(Config.key == key).first()
    if config:
        config.value = data.value
    else:
        config = Config(key=key, value=data.value)
        db.add(config)
    db.commit()
    db.refresh(config)
    return config
