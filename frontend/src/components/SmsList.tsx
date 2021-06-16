import * as React from 'react'
import useSWR from 'swr'

const fetcher = (url: string) => fetch(url).then(r => r.json())


export const SmsList = () => {
    const {data} = useSWR('http://localhost:8080/api/message', fetcher)

    console.log(data)

    return <div>stuff</div>
}