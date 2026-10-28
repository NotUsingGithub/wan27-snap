"""Upload and release build/wan27_0.1.0_{amd64,arm64}.snap to Snapcraft stable channel."""
import os, time, pathlib
from craft_store import UbuntuOneStoreClient, endpoints

client = UbuntuOneStoreClient(
    base_url="https://dashboard.snapcraft.io",
    storage_base_url="https://storage.snapcraftcontent.com",
    auth_url="https://login.ubuntu.com",
    endpoints=endpoints.U1_SNAP_STORE,
    application_name="wan27-upload",
    user_agent="wan27-upload/0.1.0",
    environment_auth="SNAPCRAFT_STORE_CREDENTIALS",
    ephemeral=True,
)

if not os.environ.get("SNAPCRAFT_STORE_CREDENTIALS"):
    kw = dict(email=os.environ["UBUNTU_ONE_EMAIL"], password=os.environ["UBUNTU_ONE_PASSWORD"])
    if os.environ.get("UBUNTU_ONE_OTP"):
        kw["otp"] = os.environ["UBUNTU_ONE_OTP"]
    client.login(
        permissions=["package_access", "package_push", "package_update", "package_release"],
        description="wan27 snap upload",
        ttl=3600,
        packages=[endpoints.Package("wan27", "snap")],
        **kw,
    )

for arch in ("amd64", "arm64"):
    path = pathlib.Path(f"build/wan27_0.1.0_{arch}.snap")
    upload_id = client.upload_file(filepath=path)
    r = client.request(
        "POST",
        "https://dashboard.snapcraft.io/dev/api/snap-push/",
        json={"name": "wan27", "updown_id": upload_id, "series": "16", "channels": ["stable"]},
    )
    status_url = r.json()["status_details_url"]
    for _ in range(60):
        s = client.request("GET", status_url).json()
        if s.get("processed"):
            print(arch, "revision:", s.get("revision"), "errors:", s.get("errors"))
            break
        time.sleep(5)
