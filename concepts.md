# ShadowNet
URL Pipeline Components:
- Downloader
- Transformer(s)
- Encryptor
- Signer?

So, in general URL will have the following structure:
[Downloader ID and parameters in base64].[TransformerID and parameters in base64].[Encryptor ID and parameters]

After decoding it from base64 all parts (at least for built-in downloaders) will have the following structure:
[ID]:[Parameter string]

## Storage
