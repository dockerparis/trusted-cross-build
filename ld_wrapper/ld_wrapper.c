#include <stdio.h>

#include <sys/types.h>
#include <sys/socket.h>
#include <stdlib.h>

#include <errno.h>

#define __USE_GNU
#include <dlfcn.h>


int execve(const char *filename, char *const argv[], char *const envp[]) {
  fprintf(stderr,"COUCOU\n");
  return 0;
}





/*int  connect(int  sockfd,  const  struct sockaddr *serv_addr, socklen_t
             addrlen){
  static int (*connect_real)(int, const  struct sockaddr*, socklen_t)=NULL;
  unsigned char *c;
  int port,ok=1;

  if (!connect_real) connect_real=dlsym(RTLD_NEXT,"connect");

  if (serv_addr->sa_family==AF_INET6) return EACCES;

  if (serv_addr->sa_family==AF_INET){
    c=serv_addr->sa_data;
    port=256*c[0]+c[1];
    c+=2;
    ok=0;

    // Allow all contacts with localhost
    if ((*c==127)&&(*(c+1)==0)&&(*(c+2)==0)&&(*(c+3)==1)) ok=1;

    // Allow contact to any WWW cache on 8080
    if (port==8080) ok=1;
  }

    if (ok) return connect_real(sockfd,serv_addr,addrlen);

    if (getenv("WRAP_TCP_DEBUG"))
      fprintf(stderr,"connect() denied to address %d.%d.%d.%d port %d\n",
              (int)(*c),(int)(*(c+1)),(int)(*(c+2)),(int)(*(c+3)),
              port);

    return EACCES;
  }
*/

  /*
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



#include <unistd.h>

int execve(const char *filename, char *const argv[], char *const envp[]) {
  write(1, "COUCOU\n", 7);
  return 0;
}
*/
