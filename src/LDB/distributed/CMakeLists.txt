FILE(GLOB dist_proto "proto/*.proto")
protobuf_generate_cpp(dist_proto_srcs dist_proto_hdrs ${dist_proto})
LIST(APPEND dist_srcs ${dist_proto_srcs})

ADD_CUSTOM_TARGET(
  copy_dist_protobuf ALL
  COMMAND ${CMAKE_COMMAND} -E make_directory "${CMAKE_BINARY_DIR}/distributed/proto"
  DEPENDS ${dist_proto_hdrs}
  COMMENT "Copying Protobuf headers"
)
FOREACH(fil ${dist_proto_hdrs})
  GET_FILENAME_COMPONENT(filename ${fil} NAME)
  ADD_CUSTOM_COMMAND(
    TARGET copy_dist_protobuf POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy_if_different ${fil}
      "${CMAKE_BINARY_DIR}/distributed/proto/${filename}"
    DEPENDS ${dist_proto}
  )
ENDFOREACH()

AUX_SOURCE_DIRECTORY(timeserver timeserver_srcs)
LIST(APPEND dist_srcs ${timeserver_srcs})

AUX_SOURCE_DIRECTORY(replication replication_srcs)
LIST(APPEND dist_srcs ${replication_srcs})

AUX_SOURCE_DIRECTORY(replication/common replication_common_srcs)
LIST(APPEND dist_srcs ${replication_common_srcs})

AUX_SOURCE_DIRECTORY(replication/vr replication_vr_srcs)
LIST(APPEND dist_srcs ${replication_vr_srcs})

AUX_SOURCE_DIRECTORY(store/common store_common_srcs)
LIST(APPEND dist_srcs ${store_common_srcs})

AUX_SOURCE_DIRECTORY(store/common/backend store_common_backend_srcs)
LIST(APPEND dist_srcs ${store_common_backend_srcs})

AUX_SOURCE_DIRECTORY(store/common/frontend store_common_frontend_srcs)
LIST(APPEND dist_srcs ${store_common_frontend_srcs})

AUX_SOURCE_DIRECTORY(store/strongstore strongstore_srcs)
LIST(APPEND dist_srcs ${strongstore_srcs})

AUX_SOURCE_DIRECTORY(lib lib_srcs)
LIST(APPEND dist_srcs ${lib_srcs})

ADD_LIBRARY(dist_objs OBJECT ${dist_srcs})
ADD_DEPENDENCIES(dist_objs copy_dist_protobuf)

if (BUILD_STATIC_LIBS)
  ADD_LIBRARY(dist STATIC $<TARGET_OBJECTS:dist_objs>)
else()
  ADD_LIBRARY(dist SHARED $<TARGET_OBJECTS:dist_objs>)
endif()
LIST(APPEND DIST_LINKER_LIBS event_openssl)
LIST(APPEND DIST_LINKER_LIBS tbb)
LIST(APPEND DIST_LINKER_LIBS crypto)
LIST(APPEND DIST_LINKER_LIBS ${PROTOBUF_LIBRARY})
TARGET_LINK_LIBRARIES(dist ${DIST_LINKER_LIBS})

ADD_EXECUTABLE(microBenchmark "exes/microBenchmark.cc")
ADD_DEPENDENCIES(microBenchmark copy_dist_protobuf)
ADD_DEPENDENCIES(microBenchmark dist)
ADD_DEPENDENCIES(microBenchmark ledger)
TARGET_LINK_LIBRARIES(microBenchmark dist ledger)
SET_TARGET_PROPERTIES(microBenchmark PROPERTIES LINK_FLAGS "${LINK_FLAGS}")

ADD_EXECUTABLE(rangeBenchmark "exes/rangeBenchmark.cc")
ADD_DEPENDENCIES(rangeBenchmark copy_dist_protobuf)
ADD_DEPENDENCIES(rangeBenchmark dist)
ADD_DEPENDENCIES(rangeBenchmark ledger)
TARGET_LINK_LIBRARIES(rangeBenchmark dist ledger)
SET_TARGET_PROPERTIES(rangeBenchmark PROPERTIES LINK_FLAGS "${LINK_FLAGS}")

ADD_EXECUTABLE(updateBenchmark "exes/updateBenchmark.cc")
ADD_DEPENDENCIES(updateBenchmark copy_dist_protobuf)
ADD_DEPENDENCIES(updateBenchmark dist)
ADD_DEPENDENCIES(updateBenchmark ledger)
TARGET_LINK_LIBRARIES(updateBenchmark dist ledger)
SET_TARGET_PROPERTIES(updateBenchmark PROPERTIES LINK_FLAGS "${LINK_FLAGS}")


ADD_EXECUTABLE(lineageBenchmark "exes/lineageBenchmark.cc")
ADD_DEPENDENCIES(lineageBenchmark copy_dist_protobuf)
ADD_DEPENDENCIES(lineageBenchmark dist)
ADD_DEPENDENCIES(lineageBenchmark ledger)
TARGET_LINK_LIBRARIES(lineageBenchmark dist ledger)
SET_TARGET_PROPERTIES(lineageBenchmark PROPERTIES LINK_FLAGS "${LINK_FLAGS}")
