import useSWR, { mutate } from 'swr'
import './App.css'
import { Box , Text} from '@chakra-ui/react';

export const ENDPOINT = "http://localhost:4000";

interface Todo {
  id : number,
  title : string,
  body : string,
  done : boolean
}

async const fetcher = (url: string) => {
  await fetch(`${ENDPOINT}/${url}`).then((res) => {
    console.log(res)
    res.json()
  });
}
function App() {
  const {data, mutate} = useSWR<Todo[]>("todo", fetcher)
  return (
    <Box>
      
    </Box>
  )
}

export default App
