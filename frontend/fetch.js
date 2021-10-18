export async function fetchWithTimeout(resource, options = {}) {
    const { timeout = 8000 } = options;
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    return new Promise((resolve, reject) => {
        fetch(resource, {
            ...options,
            signal: controller.signal  
        }).then(( response ) => {
            if (!response.ok) {
                reject(`unexpected status code: ${response.status}`)
                return
            }
            resolve(response)
        }).catch(( error ) => {
            reject(error)
        }).finally(() => {
            clearTimeout(timeoutId);
        })
    })
}