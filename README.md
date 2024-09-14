There are two api calls, GET get_text and POST store_text, both resting on a http server built with the golang stdlib. The text is stored and retrieved from a sqlite db.

A call to get_text returns `["response":[string]]` where the string slice is the list of responses.<br/>
A call to store_text requires the json payload `["text": "yourtextgoeshere"]` and the header `content-type:applications/json`. 

They also make the necessary http errors if not used properly.<br/>
The frontend has been modified to conform to the specification above.

Since the judges might not understand the Go toolchain, with ***immense generosity*** I've provided native binaries that run on arm64/darwin and amd64/linux.
They are requested to clone the repository and execute the binary which works on their system. The webpage will be active on localhost:5050.
