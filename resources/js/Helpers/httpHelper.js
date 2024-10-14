export const httpGet = async (url, params) => {
    const response = await fetch(url, {
        headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json'
        }
    })

    return response.json()
}


