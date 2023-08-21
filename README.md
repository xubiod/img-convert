# img-convert

A simple image converter, powered by Go and Dear ImGui.

Basic usage:

1. Select an export type from the list
2. Drag images to convert into that type onto the main window
3. Converts the images into the same directory as the imported file

That's the basic gist. File types with settings default to the highest settings, but you
can further change the settings under options as you please.

As of the development of internal version 0.4, there are also 6 filter options available
for mass filtering as well. There are only 6 options to keep the main vision of fast image
conversion still the main focus. Said filters are:

- Integer 1:1 upscaling
- Flip horizontally
- Flip vertically
- Rotate 90 degrees clockwise
- Rotate 180 degrees clockwise
- Rotate 270 degrees clockwise

More filters for this program are not being considered for implementation as of writing.

Written with <3 in Visual Studio Code and GoLand.
