cmake_minimum_required (VERSION 3.8)

project (whirl)

if(MSVC)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} /std:c++17")
endif(MSVC)

find_package(LLVM REQUIRED CONFIG)
message(STATUS "Found LLVM ${LLVM_PACKAGE_VERSION}")
message(STATUS "Using LLVMConfig.cmake in: ${LLVM_DIR}")

include_directories(${LLVM_INCLUDE_DIRS})
add_definitions(${LLVM_DEFINITIONS})

add_executable (whirl "src/compiler.cpp" "src/compiler.h")

llvm_map_components_to_libnames(llvm_libs support core irreader)
target_link_libraries(whirl ${llvm_libs})
