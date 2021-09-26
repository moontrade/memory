cmake_minimum_required(VERSION 3.18)
project(rpmalloc C)

set(CMAKE_C_STANDARD 11)

add_library(rpmalloc SHARED library.c rpmalloc.c)
