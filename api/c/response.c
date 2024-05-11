//
// Created by calvin on 24-5-10.
//

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "utils.c"



// Function to process reply message
void reply_message(const char* reply, char** result, Error* error) {
    if (reply == NULL || strlen(reply) < 2 || reply[strlen(reply) - 1] != 'D') {
        *result = strdup("");
        *error = REPLICATION_ERROR;
        return;
    }

    char* temp = NULL;
    *result = NULL;
    *error = NO_ERROR;

    if (reply[0] == NUMBER) {
        printf("num %s\n", reply);
        *result = bytes_to_str(reply + 1);
    } else if (reply[0] == STRING) {
        *result = bytes_to_str(reply + 1);
    } else if (reply[0] == LIST) {
        printf("list %s\n", reply);
        *result = strdup(reply + 1);
    } else {
        *result = strdup(MESSAGE_MAP[reply[0]]);
    }
}

// Class to handle response
typedef struct {
    char* reply;
    Error err;
} new_response;

// Constructor
new_response* new_response_init(const char* reply) {
    new_response* response = (new_response*)malloc(sizeof(new_response));
    if (response == NULL) {
        // Handle memory allocation failure
        return NULL;
    }

    reply_message(reply, &response->reply, &response->err);
    return response;
}

// Function to get reply
const char* get_reply(new_response* response) {
    if (response->err != NO_ERROR) {
        // Handle error case
        fprintf(stderr, "Error: %s\n", MESSAGE_MAP[response->err]);
        return NULL;
    }
    return response->reply;
}

// Function to get error
Error get_err(new_response* response) {
    return response->err;
}

// Destructor
void new_response_destroy(new_response* response) {
    if (response != NULL) {
        free(response->reply);
        free(response);
    }
}

int main() {
    const char* example_reply = "N1234567D"; // Example reply
    new_response* response = new_response_init(example_reply);
    if (response != NULL) {
        const char* reply = get_reply(response);
        if (reply != NULL) {
            printf("Reply: %s\n", reply);
        }
        Error err = get_err(response);
        if (err != NO_ERROR) {
            printf("Error: %s\n", MESSAGE_MAP[err]);
        }
        new_response_destroy(response);
    }
    return 0;
}
