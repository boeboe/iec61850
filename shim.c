#include <iec61850_client.h>
#include <mms_server.h>
#include <stdbool.h>
#include <stdint.h>

/* --- IedConnection getFile (existing) --- */
extern bool goFileHandlerCallback(void* parameter, uint8_t* buffer, uint32_t bytesRead);

void callGetFile(IedConnection conn, IedClientError* error, const char* fileName, uintptr_t handlerParam) {
	IedConnection_getFile(conn, error, fileName, goFileHandlerCallback, (void*)handlerParam);
}

/* --- MmsServer install* handlers: C shims call Go; externs declared here to avoid cgo export conflicts --- */
extern int readJournalBridgeGo(void* param, MmsDomain* domain, const char* logName, MmsServerConnection conn);
extern int getNameListBridgeGo(void* param, int nameListType, MmsDomain* domain, MmsServerConnection conn);
extern int obtainFileBridgeGo(void* param, MmsServerConnection conn, const char* sourceFilename, const char* destinationFilename);
extern void getFileCompleteBridgeGo(void* param, MmsServerConnection conn, const char* destinationFilename);

static bool readJournalHandlerShim(void* param, MmsDomain* domain, const char* logName, MmsServerConnection connection) {
	return readJournalBridgeGo(param, domain, logName, connection) != 0;
}

static bool getNameListHandlerShim(void* param, MmsGetNameListType nameListType, MmsDomain* domain, MmsServerConnection connection) {
	return getNameListBridgeGo(param, (int)nameListType, domain, connection) != 0;
}

static bool obtainFileHandlerShim(void* param, MmsServerConnection connection, const char* sourceFilename, const char* destinationFilename) {
	return obtainFileBridgeGo(param, connection, sourceFilename, destinationFilename) != 0;
}

static void getFileCompleteHandlerShim(void* param, MmsServerConnection connection, const char* destinationFilename) {
	getFileCompleteBridgeGo(param, connection, destinationFilename);
}

/* Export the shims so Go can take their address for MmsServer_install* (static would hide them). */
bool readJournalHandlerShimExport(void* param, MmsDomain* domain, const char* logName, MmsServerConnection connection) {
	return readJournalHandlerShim(param, domain, logName, connection);
}
bool getNameListHandlerShimExport(void* param, MmsGetNameListType nameListType, MmsDomain* domain, MmsServerConnection connection) {
	return getNameListHandlerShim(param, nameListType, domain, connection);
}
bool obtainFileHandlerShimExport(void* param, MmsServerConnection connection, const char* sourceFilename, const char* destinationFilename) {
	return obtainFileHandlerShim(param, connection, sourceFilename, destinationFilename);
}
void getFileCompleteHandlerShimExport(void* param, MmsServerConnection connection, const char* destinationFilename) {
	getFileCompleteHandlerShim(param, connection, destinationFilename);
}
