function clamp(value, minVal, maxVal)
    return math.max(minVal, math.min(maxVal, value))
end

function stepLighting(intensity, steps)
    return math.floor(intensity * steps) / steps
end

-- Simple distance function
function distance(p1, p2)
    local dx = p2[1] - p1[1]
    local dy = p2[2] - p1[2]
    return math.sqrt(dx * dx + dy * dy)
end

-- Simple noise function using sin/cos
function noise(x, y)
    -- Use different prime numbers and offsets to break symmetry
    local n1 = math.sin(x * 0.0234 + y * 0.0157) * math.cos(x * 0.0345 - y * 0.0287)
    local n2 = math.sin(x * 0.0567 - y * 0.0423) * math.cos(x * 0.0789 + y * 0.0634) * 0.7
    local n3 = math.sin(x * 0.1234 + y * 0.0987) * math.cos(x * 0.1456 - y * 0.1123) * 0.4
    local n4 = math.sin(x * 0.2345 - y * 0.1987) * math.cos(x * 0.2678 + y * 0.2134) * 0.2

    local combined = n1 + n2 + n3 + n4
    return combined * 0.3 + 0.5  -- normalize and reduce intensity
end

local WIDTH = 960
local HEIGHT = 540
local LIGHT_POINT = {50, 50}  -- Slightly inset from corner
local PIXEL_SIZE = 2.0

-- Base color to match your canvas (adjust these values to match your background)
local BASE_R = 87   -- Teal/green base color
local BASE_G = 131
local BASE_B = 138

-- Create new sprite
local sprite = Sprite(WIDTH, HEIGHT, ColorMode.RGB)
app.activeSprite = sprite

local layer = sprite:newLayer()
layer.name = "Vignette"

-- Create image for the layer
local image = Image(WIDTH, HEIGHT, ColorMode.RGB)

for y = 0, HEIGHT - 1 do
    for x = 0, WIDTH - 1 do
        local steppedX = math.floor(x / PIXEL_SIZE) * PIXEL_SIZE
        local steppedY = math.floor(y / PIXEL_SIZE) * PIXEL_SIZE
        local steppedPos = {steppedX, steppedY}

        -- Distance from light point
        local dis = distance(LIGHT_POINT, steppedPos)

        -- Create radial falloff with smoother curve
        local intensity = clamp(1.0 - dis / 800.0, 0.0, 1.0)
        intensity = intensity * intensity  -- Square for smoother falloff
        intensity = stepLighting(intensity, 20)

        -- Add noise
        local noiseValue = noise(steppedX, steppedY)
        intensity = intensity * (0.7 + noiseValue * 0.3)  -- Vary intensity with noise

        -- Calculate warm light color that blends with base
        local lightR = BASE_R + math.floor(60 * intensity)   -- Add warm yellow-orange
        local lightG = BASE_G + math.floor(45 * intensity)   -- Add some warmth
        local lightB = BASE_B + math.floor(10 * intensity)   -- Keep blue low for warmth

        local color = Color{
            r = clamp(lightR, 0, 255),
            g = clamp(lightG, 0, 255),
            b = clamp(lightB, 0, 255),
            a = 255
        }

        -- Set pixel
        image:drawPixel(x, y, color)
    end
end

-- Apply the image to the cel
local cel = sprite:newCel(layer, 1, image)

-- Set blend mode to normal or multiply for better blending
layer.blendMode = BlendMode.NORMAL

-- Refresh the display
app.refresh()

print("Lighting effect generated!")
print("Layer created: " .. layer.name)

-- Optional: Show dialog with parameters used
app.alert{
    title="Lighting Effect Generated",
    text={
        "Parameters used:",
        "Light Position: " .. LIGHT_POINT[1] .. ", " .. LIGHT_POINT[2],
        "Pixel Size: " .. PIXEL_SIZE,
        "",
        "You can modify the parameters at the top of the script and run again."
    }
}