import useSWR, { mutate } from 'swr'
import './App.css'
import { Box ,  List,
  ListItem,
  ListIcon,
  OrderedList,
  Heading,
  CheckboxIcon,} from '@chakra-ui/react';

export const ENDPOINT = "http://localhost:4000";

interface Todo {
  id : number,
  title : string,
  body : string,
  done : boolean
}

const fetcher = (url: string) =>
  fetch(`${ENDPOINT}/${url}`).then((r) => r.json());

function App() {
  const {data, mutate} = useSWR<Todo[]>("todo", fetcher)
  return (
    <Box>
      <Heading>
        Todos
      </Heading>
      <OrderedList>
        {data?.map((todo) => {
          return (
            <ListItem>
              {todo.done ?
              (<ListIcon as={CheckboxIcon} color='green.500' />)
              : (<ListIcon as={CheckboxIcon} color='red'/>)
              }
              
              {todo.title}
            </ListItem>
          )
        })}
      </OrderedList>
    </Box>
  )
}

export default App
