#define _GNU_SOURCE

#include <stdio.h>
#include <dlfcn.h>

static void* (*real_malloc)(size_t)=NULL;

static void mtrace_init(void)
{
  fprintf(stderr, "COUCOU\n");
  real_malloc = dlsym(RTLD_NEXT, "malloc");
  if (NULL == real_malloc) {
    fprintf(stderr, "Error in `dlsym`: %s\n", dlerror());
  }
}

void *malloc(size_t size)
{
  if(real_malloc==NULL) {
    mtrace_init();
  }

  void *p = NULL;
  fprintf(stderr, "malloc(%d) = ", (int)size);
  p = real_malloc(size);
  fprintf(stderr, "%p\n", p);
  return p;
}



/*#include <unistd.h>

int execve(const char *filename, char *const argv[], char *const envp[]) {
  write(1, "COUCOU\n", 7);
  return 0;
}
*/
