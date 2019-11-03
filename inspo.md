## Inspiration

In a natural disaster situation, infrastructure may become unavailable. This can make it difficult to find where people are, or gather information about their condition. We built a system that securely sends encrypted information with no viewable meta data about the population over a peer to peer connection. Once one peer is able to come online, the data may be extracted or posted to a central server where it can be analyzed to more appropriately distribute aid.

## What it does

Peer discovery - Currently happens over any local network or ad hoc connection, may also occur over a hotspot or other bluetooth service.

Every peer is given a unique id and tries to discover other peers. Every time a user updates their location, rescue status, or other information that may be useful for first responders a message given a timestamps and is encrypted with the ems public key and stored locally.

When a peer comes in contact with another peer it pulls all messages from that peer periodically and stores non duplicate messages in it's database.

When a peer can connect to the central server (in the case of rescue or an internet connection) all messages are sent to the server to be decrypted and analyzed.

Then first responders may analyze this information using whatever model they choose to more effectively distribute aid to those in need.

## How we built it
Messaging service - Go & Vue.js
Cloud server - Go
IBM Cloud VPS - Go & Docker
Data Visualization - Python with OpenCV, NumPy and Matplotlib)

## Challenges we ran into
Peer-to-peer connection
Reliable geolocation updates
Data store 
ML planning wasn't implemented

## Accomplishments that we're proud of
The potential to help victims of natural disasters.
Implemented a lot in a short amount of time.
Fleshed out idea -- we thought about and discussed a lot of ways to build off of this MVP.
Proof of concept for mobile device Bluetooth communication.

## What we learned
- Go Storm DB
- Writing an API in Go
- Vue.js
- Zeroconf and service discovery

## What's next for Disasterpeer
- Finish Mobile application
- Data model refactor
- Support bluetooth 
- Gathering additional information
    - Who each peer came in contact with
    - Extra meta field
- Decide where to send aid to most effectively help natural disaster victims                    

