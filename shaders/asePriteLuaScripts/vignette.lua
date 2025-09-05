

function clamp(value, minVal, maxVal)
    return math.max(minVal, math.min(maxVal, value))
end
function stepLighting(intensity, steps)
    return math.floor(intensity * steps) / steps
end
function smoothstep(edge0, edge1, x)
    local t = clamp((x - edge0) / (edge1 - edge0), 0.0, 1.0)
    return t * t * (3.0 - 2.0 * t)
end

local WIDTH = 960
local HEIGHT = 540
local VIGNETTE_STRENGTH = 0.9
local VIGNETTE_SIZE = 0.5
local CENTER_X = 0.5
local CENTER_Y = 0.5
local FEATHER = 0.7

-- Create new sprite
local sprite = Sprite(WIDTH, HEIGHT, ColorMode.RGB)
app.activeSprite = sprite

-- Create a new layer for the vignette effect
local layer = sprite:newLayer()
layer.name = "Vignette"

-- Create image for the layer
local image = Image(WIDTH, HEIGHT, ColorMode.RGB)

-- Calculate center position in pixels
local centerX = WIDTH * CENTER_X
local centerY = HEIGHT * CENTER_Y

-- Calculate maximum distance from center to corner
local maxDistance = math.sqrt(
        math.max((centerX)^2, (WIDTH - centerX)^2) +
                math.max((centerY)^2, (HEIGHT - centerY)^2)
)

-- Generate vignette effect
for y = 0, HEIGHT - 1 do
    for x = 0, WIDTH - 1 do
        -- Calculate distance from center
        local dx = x - centerX
        local dy = y - centerY
        local distance = math.sqrt(dx * dx + dy * dy)

        -- Normalize distance (0.0 at center, 1.0 at max distance)
        local normalizedDistance = distance / maxDistance

        -- Create vignette falloff
        local vignetteRadius = VIGNETTE_SIZE
        local vignetteFalloff = stepLighting(vignetteRadius, vignetteRadius + (1.0 - vignetteRadius) / FEATHER, normalizedDistance)

        -- Apply vignette strength
        local vignetteAmount = vignetteFalloff * VIGNETTE_STRENGTH

        -- Create the darkening effect (black with alpha based on vignette amount)
        -- For a very light vignette, we'll use a subtle darkening
        local darkness = math.floor(vignetteAmount * 255)

        local color = Color{
            r = 0,
            g = 0,
            b = 0,
            a = darkness  -- Alpha controls how much darkening
        }

        -- Set pixel
        image:drawPixel(x, y, color)
    end
end

-- Apply the image to the cel
local cel = sprite:newCel(layer, 1, image)

-- Set blend mode to multiply for darkening effect
layer.blendMode = BlendMode.MULTIPLY
layer.opacity = 255  -- You can lower this for an even lighter effect

-- Refresh the display
app.refresh()

print("Vignette effect generated!")
print("Layer created: " .. layer.name)
print("Blend mode: Multiply")
print("You can adjust the layer opacity for a lighter effect.")

-- Show dialog with parameters used
app.alert{
    title="Vignette Effect Generated",
    text={
        "Parameters used:",
        "Vignette Strength: " .. VIGNETTE_STRENGTH,
        "Vignette Size: " .. VIGNETTE_SIZE,
        "Center: (" .. CENTER_X .. ", " .. CENTER_Y .. ")",
        "Feather: " .. FEATHER,
        "",
        "Adjust layer opacity for lighter effect.",
        "Use Multiply blend mode for darkening."
    }
}