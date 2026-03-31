"""Package management API."""
from typing import List
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from ..database import get_db
from ..schemas import PackageCreate, PackageResponse, PackageDetail
from ..services import PackageService

router = APIRouter()


@router.get("/", response_model=List[PackageResponse])
def list_packages(db: Session = Depends(get_db)):
    """List all packages."""
    service = PackageService(db)
    return service.list_packages()


@router.get("/{package_id}", response_model=PackageDetail)
def get_package(package_id: int, db: Session = Depends(get_db)):
    """Get package details."""
    service = PackageService(db)
    package = service.get_package(package_id)
    if not package:
        raise HTTPException(status_code=404, detail="Package not found")
    return package


@router.post("/", response_model=PackageResponse)
def import_package(data: PackageCreate, db: Session = Depends(get_db)):
    """Import a package."""
    service = PackageService(db)
    try:
        if data.source_type == "git":
            return service.import_from_git(data.url_or_path, data.name)
        elif data.source_type == "local":
            return service.import_from_local(data.url_or_path, data.name)
        else:
            raise HTTPException(status_code=400, detail=f"Unsupported source type: {data.source_type}")
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.delete("/{package_id}")
def delete_package(package_id: int, db: Session = Depends(get_db)):
    """Delete a package."""
    service = PackageService(db)
    if not service.delete_package(package_id):
        raise HTTPException(status_code=404, detail="Package not found")
    return {"status": "deleted"}
