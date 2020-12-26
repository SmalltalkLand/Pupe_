cmake_minimum_required(VERSION 3.10)

# set the project name
project(Pupe)

# add the executable
add_executable(Pupe Pupe.cxx)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

add_subdirectory(coroutine)

target_link_libraries(Pupe PUBLIC coroutine)

target_include_directories(Pupe PUBLIC
                          "${PROJECT_SOURCE_DIR}/coroutine/interface"
                          )
FIND_PACKAGE(PkgConfig REQUIRED)
PKG_CHECK_MODULES(GTK3 REQUIRED gtk+-3.0)

# Setup CMake to use GTK+, tell the compiler where to look for headers
# and to the linker where to look for libraries
INCLUDE_DIRECTORIES(${GTK3_INCLUDE_DIRS})
LINK_DIRECTORIES(${GTK3_LIBRARY_DIRS})

# Add other flags to the compiler
ADD_DEFINITIONS(${GTK3_CFLAGS_OTHER})

# Add an executable compiled from hello.c

# Link the target to the GTK+ libraries
TARGET_LINK_LIBRARIES(Pupe ${GTK3_LIBRARIES})