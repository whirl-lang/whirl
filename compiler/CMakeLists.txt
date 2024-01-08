cmake_minimum_required (VERSION 3.8)

project (llvm)

find_package(LLVM REQUIRED CONFIG)
include_directories(${LLVM_INCLUDE_DIRS})
add_definitions(${LLVM_DEFINITIONS})

add_executable (llvm "llvm/llvm.cpp" "llvm/llvm.h")

llvm_map_components_to_libnames(llvm_libs support core irreader)
target_link_libraries(llvm ${llvm_libs})
