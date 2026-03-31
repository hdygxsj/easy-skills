"""Installation management API."""
from typing import List
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from ..database import get_db
from ..schemas import InstallationCreate, InstallationResponse
from ..services import InstallService

router = APIRouter()


@router.get("/", response_model=List[InstallationResponse])
def list_installations(db: Session = Depends(get_db)):
    """List all installations."""
    service = InstallService(db)
    return service.list_installations()


@router.get("/{installation_id}", response_model=InstallationResponse)
def get_installation(installation_id: int, db: Session = Depends(get_db)):
    """Get installation details."""
    service = InstallService(db)
    installation = service.get_installation(installation_id)
    if not installation:
        raise HTTPException(status_code=404, detail="Installation not found")
    return installation


@router.post("/", response_model=InstallationResponse)
def install_group(data: InstallationCreate, db: Session = Depends(get_db)):
    """Install a group to IDE."""
    service = InstallService(db)
    try:
        return service.install_group(
            data.group_id,
            data.target_ide,
            data.install_scope,
            data.install_path
        )
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))


@router.delete("/{installation_id}")
def uninstall(installation_id: int, db: Session = Depends(get_db)):
    """Uninstall by ID."""
    service = InstallService(db)
    if not service.uninstall(installation_id):
        raise HTTPException(status_code=404, detail="Installation not found")
    return {"status": "uninstalled"}


@router.get("/{installation_id}/compare")
def compare_versions(installation_id: int, db: Session = Depends(get_db)):
    """Compare installed vs local versions."""
    service = InstallService(db)
    result = service.compare_versions(installation_id)
    if not result:
        raise HTTPException(status_code=404, detail="Installation not found")
    return result


@router.post("/{installation_id}/upgrade", response_model=InstallationResponse)
def upgrade_installation(installation_id: int, db: Session = Depends(get_db)):
    """Upgrade installation to latest local version."""
    service = InstallService(db)
    installation = service.get_installation(installation_id)
    if not installation:
        raise HTTPException(status_code=404, detail="Installation not found")
    
    # Reinstall (switch to same group)
    try:
        return service.switch_group(
            installation_id,
            installation.group_id,
            installation.target_ide,
            installation.install_scope,
            installation.install_path if installation.install_scope == "project" else None
        )
    except Exception as e:
        raise HTTPException(status_code=400, detail=str(e))
