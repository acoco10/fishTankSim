# Fish Fish Fish!

### A Small Browser Tamagachi Esque FishScene Written in Go/ebiten compiled to WASM that can save user State 
## Goals 
- learn javascript asynch functions for loading and saving player State
- Develop a clean front-end in JavaScript and HTML
- user/server management in Go backend
- Develop good pixel art visuals for water and particles
- develop behaviour tree for Fish ui with good looking swimming movement
- automate sprite loading based on Fish type and level 

## Dev Notes
- Fish swimming behavior is pretty good; they are randomly assigned a target point that they swim towards unless there is food
- Animation looks good, Fish swim at different speeds, and it looks smooth. 
- The initial challenge of loading a compiled WASM that triggers a JavaScript file to load the save and wait for it to be completed is solved.
- In other words, you can feed the Fish, see them grow, then reload, and they will still be grown 
- For now, it just runs on localhost; the next step will be hosting/ user management. 
