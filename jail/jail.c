#include <sys/ptrace.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>
#include <sys/user.h>
#include <stdlib.h>
#include <string.h>

void child_process(char *path) {
    size_t len = strlen(path);
    char *buf = malloc(len);
    memcpy(buf, path, len);
    if(execlp(buf, buf, NULL) == -1){
        printf("failed to execve\n");
    }
}

int main(int argc, char **argv) {
    if(argc < 2) {
        printf("at least 2 args");
        return -1;
    }

    pid_t child;
    int status;
    char *child_process_name = argv[0];

    struct user_regs_struct regs;

    child = fork();
    if(child < 0) {
        printf("fork failed");
        return -1;
    }

    if(child == 0) {
        ptrace(PTRACE_TRACEME);
        kill(getpid(), 19);
        child_process(child_process_name);
    }else{
        wait(&status);
        while(1){
            ptrace(PTRACE_SYSCALL, child, NULL, NULL);
            waitpid(child, &status, NULL);
            ptrace(PTRACE_GETREGS, child, NULL, &regs);
            

            if(regs.orig_rax == 231) {
                kill(child, 9);
                break;
            }else if(regs.orig_rax == 59) {
                char *buf = malloc(1024);
                int buf_index = 0;
                int copy_processing = 1;

                do{
                    int data = ptrace(PTRACE_PEEKTEXT, child, regs.rdi + buf_index * 4, NULL);
                    if(data == -1){
                        goto end;
                    }

                    memcpy(buf + buf_index * 4, &data, 4);
                    for(int i = 0; i < 4; i++){
                        if(*(buf + buf_index * 4 + i) == 0){
                            copy_processing = 0;
                            break;
                        }
                    }
                    buf_index += 1;

                } while(copy_processing && buf_index < 255);


                if(strcmp(buf, child_process_name) != 0){
                    kill(child, 9);
                    return -2;
                }else{
                    printf("EXECVE: %s\n", buf);
                }
                end:
                    printf("allowed execve\n");
            }else if(regs.orig_rax == 2){
                kill(child, 9);
                break;
            }else{
                printf("SYSCALL %lld\n", regs.orig_rax);
            }
        }
    }
}