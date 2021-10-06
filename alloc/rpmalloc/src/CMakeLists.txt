cmake_minimum_required(VERSION 3.18)
project(rpmalloc C)

set(CMAKE_C_STANDARD 11)

if(UNIX AND NOT APPLE)
    set(LINUX TRUE)
endif()

if(APPLE)
	set(CMAKE_C_FLAGS "-Ofast -fPIC")
endif()
if(LINUX)
    set(CMAKE_C_FLAGS "-Ofast -fPIC -D_GNU_SOURCE")
endif()

add_library(rpmalloc SHARED rpmalloc.c)
