from fastapi import FastAPI, UploadFile, File
import uvicorn
import random
import time
import uuid

app = FastAPI(title="Smart Factory AI QC Service (In-Memory)")

# In-Memory Database Simulators
db_smartfactory_qc = []  # List of inspection results
minio_storage = {}       # Dictionary mapping URL/Path -> File Info

def dummy_inference_resnet18(filename: str) -> dict:
    """Mensimulasikan hasil deteksi model AI (ResNet18)"""
    # 1. Simulasi waktu inferensi (misal model butuh waktu untuk memproses)
    time.sleep(random.uniform(0.1, 0.5))
    
    # 2. Acak status dengan probabilitas tertentu (misal 80% Normal, 20% Defect)
    status = "Normal" if random.random() < 0.8 else "Defect"
    
    # 3. Confidence score antara 0.85 hingga 0.99
    confidence = round(random.uniform(0.85, 0.99), 4)
    
    print(f"[Model ResNet18] Menganalisis gambar: {filename}")
    print(f" -> Hasil Prediksi: {status} (Confidence: {confidence*100:.2f}%)")
    
    return {
        "status": status,
        "confidence": confidence
    }

@app.post("/api/v1/qc/inspect")
async def inspect_product(image: UploadFile = File(...)):
    print(f"\n[Kamera Industri] 📸 Mengunggah gambar baru: {image.filename}")
    
    # 1. Simpan ke MinIO Simulator (In-Memory)
    file_id = str(uuid.uuid4())
    virtual_url = f"http://minio:9000/smartfactory-images/{file_id}_{image.filename}"
    minio_storage[virtual_url] = {
        "filename": image.filename,
        "content_type": image.content_type,
        "uploaded_at": time.time()
    }
    print(f"[MinIO] Gambar disimpan di path virtual: {virtual_url}")
    
    # 2. Jalankan Inferensi AI
    inference_result = dummy_inference_resnet18(image.filename)
    
    # 3. Simpan Hasil ke db_smartfactory_qc Simulator (In-Memory)
    inspection_record = {
        "id": str(uuid.uuid4()),
        "image_url": virtual_url,
        "decision": inference_result["status"],
        "confidence_score": inference_result["confidence"],
        "inspected_at": time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    }
    db_smartfactory_qc.append(inspection_record)
    print(f"[PostgreSQL QC] Rekam inspeksi disimpan dengan ID: {inspection_record['id']}")
    
    return {
        "message": "Inspeksi berhasil",
        "result": inspection_record
    }

if __name__ == "__main__":
    print("="*60)
    print("Memulai AI Quality Control Service (Pure Local Mode)...")
    print("="*60)
    uvicorn.run(app, host="0.0.0.0", port=8000)
