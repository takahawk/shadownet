# ShadowNet
URL Pipeline Components:
- Downloader
- Transformer(s)
- Encryptor
- Signer?

So, in general URL will have the following structure:
[Downloader ID and parameters in base64].([TransformerID/EncryptorID and parameters in base64])*

After decoding it from base64 all parts (at least for built-in downloaders) will have the following structure:
[Type]_[ID]:[Base64dCommaSeparatedParameters]
Type can be "enc", "trans" or "down".
## Storage

## Ideas:
Editable storage
Keyring?
How to work with OAuth Flow?

### Applications
Journal
Chat
Photo Storage
Music Streaming

### Storages
Google Drive
