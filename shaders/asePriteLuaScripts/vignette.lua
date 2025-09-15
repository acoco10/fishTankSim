-- Aseprite Caustic Generator
-- Save this as caustic_generator.lua in your Aseprite scripts folder
-- Run from Aseprite: File > Scripts > caustic_generator

-- Configuration
local config = {
    width = 128,
    height = 128,
    frames = 16,

    -- Caustic parameters
    num_waves = 6,
    wave_speed = 0.1,
    intensity = 1.2,

    -- Dithering parameters
    dither_strength = 0.3,
    threshold_levels = 4,

    -- Colors (RGB 0-255)
    bg_color = {0, 20, 40},      -- Dark blue background
    caustic_color = {180, 220, 255}, -- Light blue caustics
}

-- Bayer dithering matrix (4x4)
local bayer_matrix = {
    { 0,  8,  2, 10},
    {12,  4, 14,  6},
    { 3, 11,  1,  9},
    {15,  7, 13,  5}
}

-- Normalize bayer matrix
for i = 1, 4 do
    for j = 1, 4 do
        bayer_matrix[i][j] = bayer_matrix[i][j] / 16.0
    end
end

-- Simple noise function
function noise(x, y, seed)
    local n = math.sin(x * 12.9898 + y * 78.233 + seed) * 43758.5453
    return (n - math.floor(n))
end

-- Fractal noise
function fractal_noise(x, y, octaves, seed)
    local value = 0
    local amplitude = 1
    local frequency = 1
    local max_value = 0

    for i = 1, octaves do
        value = value + noise(x * frequency, y * frequency, seed) * amplitude
        max_value = max_value + amplitude
        amplitude = amplitude * 0.5
        frequency = frequency * 2
    end

    return value / max_value
end

-- Generate caustic pattern
function generate_caustic(x, y, time)
    local caustic = 0

    -- Multiple wave interference
    for i = 1, config.num_waves do
        local angle = (i / config.num_waves) * math.pi * 2
        local wave_x = math.cos(angle)
        local wave_y = math.sin(angle)

        -- Wave parameters
        local freq = 0.8 + i * 0.3
        local phase = time * config.wave_speed * (0.5 + i * 0.2)
        local offset = i * 123.456

        -- Calculate wave
        local dist = (x * wave_x + y * wave_y) * freq + phase + offset
        local wave = math.sin(dist * math.pi * 2)

        -- Add distortion
        local distort_x = x + noise(x * 0.1, y * 0.1, offset) * 0.2
        local distort_y = y + noise(x * 0.1 + 100, y * 0.1 + 100, offset) * 0.2
        local distort = fractal_noise(distort_x * 2, distort_y * 2, 3, offset)

        wave = wave + distort * 0.3
        caustic = caustic + wave * wave
    end

    caustic = caustic / config.num_waves
    caustic = math.max(0, caustic)
    caustic = math.pow(caustic, 0.8) * config.intensity

    -- Add turbulence
    local turbulence = fractal_noise(x * 3 + time * 0.05, y * 3 + time * 0.07, 4, 999)
    caustic = caustic + turbulence * 0.2

    return math.max(0, math.min(1, caustic))
end

-- Apply dithering
function apply_dithering(value, x, y)
    local threshold = bayer_matrix[(y % 4) + 1][(x % 4) + 1]
    threshold = threshold * config.dither_strength

    local quantized = math.floor(value * config.threshold_levels + threshold) / config.threshold_levels
    return math.max(0, math.min(1, quantized))
end

-- Mix colors
function mix_color(bg, fg, alpha)
    return {
        math.floor(bg[1] * (1 - alpha) + fg[1] * alpha),
        math.floor(bg[2] * (1 - alpha) + fg[2] * alpha),
        math.floor(bg[3] * (1 - alpha) + fg[3] * alpha)
    }
end

-- Main generation function
function generate_caustic_animation()
    -- Create new sprite
    local sprite = Sprite(config.width, config.height, ColorMode.RGB)
    sprite.filename = "caustic_animation"

    -- Generate frames
    app.transaction(function()
        for frame_num = 1, config.frames do
            local time = (frame_num - 1) / config.frames

            -- Create new frame if not first
            if frame_num > 1 then
                sprite:newFrame()
            end

            -- Get current frame
            local frame = sprite.frames[frame_num]
            frame.duration = 0.1 -- 100ms per frame

            -- Create image for this frame
            local image = Image(config.width, config.height, ColorMode.RGB)

            -- Generate caustic pattern
            for y = 0, config.height - 1 do
                for x = 0, config.width - 1 do
                    -- Normalize coordinates
                    local nx = x / config.width
                    local ny = y / config.height

                    -- Generate caustic value
                    local caustic = generate_caustic(nx, ny, time)

                    -- Apply dithering
                    local dithered = apply_dithering(caustic, x, y)

                    -- Mix background and caustic colors
                    local final_color = mix_color(config.bg_color, config.caustic_color, dithered)

                    -- Set pixel
                    local color = Color{
                        r = final_color[1],
                        g = final_color[2],
                        b = final_color[3],
                        a = 255
                    }
                    image:drawPixel(x, y, color)
                end
            end

            -- Apply image to cel
            local cel = sprite:newCel(sprite.layers[1], frame, image)

            -- Update progress
            app.refresh()
        end
    end)

    -- Set animation to loop
    sprite.frames[1].duration = 0.1

    app.alert("Caustic animation generated! " .. config.frames .. " frames created.")
end

-- Create dialog for parameters
function show_dialog()
    local dlg = Dialog("Caustic Generator")

    dlg:number{ id="width", label="Width:", text=tostring(config.width) }
    dlg:number{ id="height", label="Height:", text=tostring(config.height) }
    dlg:number{ id="frames", label="Frames:", text=tostring(config.frames) }
    dlg:separator()
    dlg:slider{ id="waves", label="Wave Count:", min=3, max=12, value=config.num_waves }
    dlg:slider{ id="speed", label="Speed:", min=0.05, max=0.3, value=config.wave_speed }
    dlg:separator()
    dlg:slider{ id="dither", label="Dither Strength:", min=0.1, max=0.8, value=config.dither_strength }
    dlg:slider{ id="levels", label="Threshold Levels:", min=2, max=8, value=config.threshold_levels }
    dlg:separator()
    dlg:button{ id="generate", text="Generate Caustics" }
    dlg:button{ id="cancel", text="Cancel" }

    dlg:show()

    if dlg.data.generate then
        -- Update config with dialog values
        config.width = dlg.data.width
        config.height = dlg.data.height
        config.frames = dlg.data.frames
        config.num_waves = dlg.data.waves
        config.wave_speed = dlg.data.speed
        config.dither_strength = dlg.data.dither
        config.threshold_levels = dlg.data.levels

        -- Generate the animation
        generate_caustic_animation()
    end
end

-- Run the dialog
show_dialog()