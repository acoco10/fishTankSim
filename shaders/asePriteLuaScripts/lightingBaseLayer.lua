-- Aseprite Lighting Effect Generator
-- Based on your shader code
-- File > Scripts > Open Script File, then select this .lua file

-- Simple config file parser (key=value format)

function stepLighting(intensity, steps)
    return math.floor(intensity * steps) / steps
end

function checkRect(x, y, tankRect)
    local xCheck = x >= tankRect[1] and x <= tankRect[3]
    local yCheck = y >= tankRect[2] and y <= tankRect[4]
    return xCheck and yCheck
end

function clamp(value, minVal, maxVal)
    return math.max(minVal, math.min(maxVal, value))
end

function distance(p1, p2)
    local dx = p1[1] - p2[1]
    local dy = p1[2] - p2[2]
    return math.sqrt(dx * dx + dy * dy)
end

function normalize(vector)
    local length = math.sqrt(vector[1] * vector[1] + vector[2] * vector[2])
    if length == 0 then
        return {0, 0}
    end
    return {vector[1] / length, vector[2] / length}
end


local WIDTH = 960
local HEIGHT = 540
local LIGHT_POINT = {197+194, 188-52}  -- {x, y}
local LIGHT_WIDTH = 120
local TANK_RECT = {197, 188, 581, 373}  -- {x1, y1, x2, y2}
local PIXEL_SIZE = 2.0


-- Create new sprite
local sprite = Sprite(WIDTH, HEIGHT, ColorMode.RGB)
app.activeSprite = sprite

-- Create a new layer for the lighting effect
local layer = sprite:newLayer()
layer.name = "Lighting Effect"



-- Create image for the layer
local image = Image(WIDTH, HEIGHT, ColorMode.RGB)

-- Generate lighting effect
for y = 0, HEIGHT - 1 do
    for x = 0, WIDTH - 1 do
        local color = Color{r=0, g=0, b=0, a=0}  -- Default black

        -- Check if inside tank
        if checkRect(x, y, TANK_RECT) then
            -- Quantize position for stepped lighting
            local steppedX = math.floor(x / PIXEL_SIZE) * PIXEL_SIZE
            local steppedY = math.floor(y / PIXEL_SIZE) * PIXEL_SIZE
            local steppedPos = {steppedX, steppedY}

            -- Calculate lighting from stepped position
            local lineHalfWidth = LIGHT_WIDTH * 0.5
            local clampedX = clamp(steppedPos[1],
                    LIGHT_POINT[1] - lineHalfWidth,
                    LIGHT_POINT[1] + lineHalfWidth)
            local closestPointOnLine = {clampedX, LIGHT_POINT[2]}

            -- Distance and direction calculations
            local dis = distance(closestPointOnLine, steppedPos)

            local lightDir = {0, 0}
            if dis > 0 then
                lightDir = normalize({steppedPos[1] - closestPointOnLine[1],
                                      steppedPos[2] - closestPointOnLine[2]})
            end

            local downwardBias = math.max(0.0, lightDir[2]) * 0.9

            -- Light intensity calculation (updated to match shader)
            local lightRadius = clamp(1.0 - dis / 250.0, 0.0, 1.0) * (1.0 + downwardBias)
            lightRadius = stepLighting(lightRadius, 15)

            -- Apply shader-style lighting (darken base, add light)
            local baseR = math.floor(255 * 0.4)  -- Darken background
            local baseG = math.floor(255 * 0.4)
            local baseB = math.floor(255 * 0.4)

            -- Add the light contribution
            local lightR = math.floor(255 * 0.45 * lightRadius)
            local lightG = math.floor(255 * 0.5 * lightRadius)
            local lightB = math.floor(255 * 0.55 * lightRadius)

            color = Color{
                r = clamp((baseR + lightR) * 0.6, 0, 255),
                g = clamp((baseG + lightG) *0.8, 0, 255),
                b = clamp(baseB + lightB, 0, 255),
                a = 255
            }

            -- For colored light, try these instead:
            -- Warm light:
            -- color = Color{
            --     r = math.floor(intensity),
            --     g = math.floor(intensity * 0.8),
            --     b = math.floor(intensity * 0.4),
            --     a = 255
            -- }

            -- Cool light:
            -- color = Color{
            --     r = math.floor(intensity * 0.4),
            --     g = math.floor(intensity * 0.6),
            --     b = math.floor(intensity),
            --     a = 255
            -- }
        end

        -- Set pixel
        image:drawPixel(x, y, color)
    end
end

-- Apply the image to the cel
local cel = sprite:newCel(layer, 1, image)

-- Set blend mode to additive (Screen mode is closest to additive in Aseprite)
layer.blendMode = BlendMode.SCREEN

-- Refresh the display
app.refresh()

print("Lighting effect generated!")
print("Layer created: " .. layer.name)
print("You can adjust the blend mode or layer opacity as needed.")

-- Optional: Show dialog with parameters used
app.alert{
    title="Lighting Effect Generated",
    text={
        "Parameters used:",
        "Light Position: (" .. LIGHT_POINT[1] .. ", " .. LIGHT_POINT[2] .. ")",
        "Light Width: " .. LIGHT_WIDTH,
        "Tank Rect: (" .. TANK_RECT[1] .. ", " .. TANK_RECT[2] .. ", " .. TANK_RECT[3] .. ", " .. TANK_RECT[4] .. ")",
        "Pixel Size: " .. PIXEL_SIZE,
        "",
        "You can modify the parameters at the top of the script and run again."
    }
}