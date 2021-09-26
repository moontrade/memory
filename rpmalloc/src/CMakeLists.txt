cmake_minimum_required(VERSION 3.18)
project(rpmalloc C)

set(CMAKE_C_STANDARD 11)

if(UNIX AND NOT APPLE)
    set(LINUX TRUE)
endif()

set(CMAKE_C_FLAGS "-O3")
if(LINUX)
    set(CMAKE_C_FLAGS "-O3 -D_GNU_SOURCE")
endif()

add_library(rpmalloc SHARED library.c rpmalloc.c)
