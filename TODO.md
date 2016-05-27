# TODO

- parallelize ch1/e04, as several gorutines can open and read files in parallel.
  The map where results are stored will need to be accessed synchronously, but
  that is already solved with a mutex.
