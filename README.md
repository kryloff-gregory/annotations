# Annotations service

## What was done

The service itself is implemented. 
As a storage LevelDB was used. 
I intended to use NoSQL KV storage as there are no queries requiring indices other that ID 
(as key for annotation is a compound videoURL:annotationID). 
Later I regretted to choose this storage, as I discovered that there is no opportunity to make pessimistic nor pessimistic locks provided. 
But it was already too late.
In order to make the calls to service consistent, a component providing locks was implemented. 
In production the locks would be acquired not in RAM, we could use Redis for it. 
<br>
For authentication/authorization JWT tokens were chosen.

## What was not done because of lack of time

Proper testing unit, integration and manual testing. (Only end-to-end DAO and minor unit tests were implemented).

Docker file.

New user registration endpoint (existing users with their hashed passwords are stored as JSON in a file).

Proper API documentation. 

Proper error processing (HTTP codes might be not very accurate nor informative).

## Example requests

### Login

curl -X POST localhost:5001/api/v1/login -d "{\"userName\":\"Bob\", \"password\":\"<your password>\"}"

### New video

curl --header "Authorization: <your token>" -X PUT localhost:5001/api/v1/video -d "{\"url\":\"https://youtu.be/AZAZA\", \"duration\":\"3m\"}"

### Get video 

curl --header "Authorization: <your token>"  localhost:5001/api/v1/video?video_url=https://youtu.be/AZAZAZ

### Delete video

curl --header "Authorization: <your token>" -X DELETE localhost:5001/api/v1/video?video_url=https://youtu.be/AZAZAZ

### Get annotations for video

curl --header "Authorization: <your token>"  localhost:5001/api/v1/annotation?video_url=https://youtu.be/AZAZAZ

### Delete annotation

curl --header "Authorization: <your token>" -X DELETE localhost:5001/api/v1/annotation?video_url=https://youtu.be/AZAZAZ&id=<annotation id>

### New annotation

curl --header "Authorization: <your token>" -X PUT localhost:5001/api/v1/annotation -d "{<check annotationAddRequest>}"

### Update annotation

curl --header "Authorization: <your token>" -X POST localhost:5001/api/v1/annotation -d "{<check annotationUpdateRequest>}"

## Assumptions

User can delete video only if there are no annotations for this video.

User can delete or update only annotations added by himself.

As a consequence, before deleting his video, user should first ask other users to delete their annotations for the video.

Video URLs should match regexp "https://youtu.be/[a-zA-Z0-9]+$"