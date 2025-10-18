# Umbrella Corp. INAN

© Umbrella Corp. Internal Network Access Node. All Rights Reserved.

### Usage (MacOS)

```sh
# Configure
/Users/<user>/Qt/Tools/CMake/CMake.app/Contents/bin/cmake -S . -B build -G Ninja -DCMAKE_PREFIX_PATH:PATH=/Users/<user>/Qt/6.9.3/macos -DCMAKE_BUILD_TYPE:STRING=Build

# Build
/Users/<user>/Qt/Tools/CMake/CMake.app/Contents/bin/cmake --build build --target all

# Run
./build/umbrella.app/Contents/MacOS/umbrella
```
