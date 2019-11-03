# Disasterpeer
Peer-to-peer encrypted messaging and cloud service designed to provide insight to first responders during a natural disaster.

## About
Nodes keep an encrypted log of messages, with location and safety status. When two devices can communicate (either via hotspot, or for mobile via Bluetooth), those devices share all messages.

When a node is able to access the internet, that node's message log is pushed to an IBM cloud instance.

## Usage

```
git clone git@github.com:jlyon1/disasterpeer.git
cd /disasterpeer
```

Run messaging application:
```
go run *.go
```

Run cloud server:
```
go run server/*.go
```
