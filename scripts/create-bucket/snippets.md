# Bucketの作成

```bash
gcloud storage buckets create gs://YOUR_TFSTATE_BUCKET \
  --project=YOUR_PROJECT \
  --location=asia-northeast1
```

# Bucketの設定

```bash
gcloud storage buckets update gs://YOUR_TFSTATE_BUCKET \
  --versioning \
  --uniform-bucket-level-access \
  --public-access-prevention \
  --lifecycle-file=lifecycle.json
```

# IAMの設定

```bash
gcloud storage buckets add-iam-policy-binding gs://YOUR_TFSTATE_BUCKET \
  --member="user:your-name@example.com" \
  --role="roles/storage.objectAdmin"
```
