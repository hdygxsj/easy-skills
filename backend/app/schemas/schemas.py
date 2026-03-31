"""Pydantic models for request/response validation."""
from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel


# Source schemas
class SourceCreate(BaseModel):
    type: str  # 'git' | 'guide' | 'local'
    url_or_path: str


class SourceResponse(BaseModel):
    id: int
    type: str
    url_or_path: str
    created_at: datetime

    class Config:
        from_attributes = True


# Package schemas
class PackageCreate(BaseModel):
    source_type: str  # 'git' | 'guide' | 'local'
    url_or_path: str
    name: Optional[str] = None


class PackageResponse(BaseModel):
    id: int
    name: str
    version: Optional[str]
    local_path: str
    created_at: datetime
    updated_at: datetime

    class Config:
        from_attributes = True


class PackageDetail(PackageResponse):
    source: Optional[SourceResponse]
    metadata_json: Optional[str]


# Group schemas
class GroupCreate(BaseModel):
    name: str
    description: Optional[str] = None


class GroupResponse(BaseModel):
    id: int
    name: str
    description: Optional[str]
    created_at: datetime

    class Config:
        from_attributes = True


class GroupDetail(GroupResponse):
    packages: List[PackageResponse]


# Installation schemas
class InstallationCreate(BaseModel):
    group_id: int
    target_ide: str  # 'qoder' | 'cursor'
    install_scope: str  # 'user' | 'project'
    install_path: Optional[str] = None  # Required for 'project' scope


class InstallationResponse(BaseModel):
    id: int
    group_id: int
    target_ide: str
    install_scope: str
    install_path: str
    installed_version: Optional[str]
    installed_at: datetime

    class Config:
        from_attributes = True


# Config schemas
class ConfigUpdate(BaseModel):
    key: str
    value: str


class ConfigResponse(BaseModel):
    key: str
    value: Optional[str]
    updated_at: Optional[datetime]

    class Config:
        from_attributes = True


# Chat schemas
class ChatMessage(BaseModel):
    message: str


class ChatResponse(BaseModel):
    response: str
    tool_calls: Optional[List[dict]] = None
