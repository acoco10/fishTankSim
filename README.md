# Fish Fish Fish!

You can try it on my site here: https://collisionposition.netlify.app/gamepage

### A Small Browser Tamagachi Esque FishScene Written in Go/ebiten compiled to WASM that can save user State 
## Goals 
- learn javascript asynch functions for loading and saving player State
- Develop a clean front-end in JavaScript and HTML
- user/server management in Go backend
- Develop good pixel art visuals for water and particles
- Develop a behaviour tree for Fish ui with a good-looking swimming movement and behaviour that is interesting to observe 
- Automate sprit loading based on Fish type and level
- native JavaScript music player/streamer that communicates with web assembly application (storing large music file within web assembly is cumbersome and cuases performance issues)

## Dev Notes
- Fish swimming behavior is pretty good; they are randomly assigned a target point that they swim towards unless there is food
- Animation looks good, Fish swim at different speeds, and it looks smooth. 
- The initial challenge of loading a compiled WASM that triggers a JavaScript file to load the save and wait for it to be completed is solved.
- In other words, you can feed the Fish, see them grow, then reload, and they will still be grown 
~~For now, it just runs on localhost; the next step will be hosting/ user management~~
- uploaded to the site and users can log in to save their game
- hosted on heroku (back end) and netlify(static/frontend)
