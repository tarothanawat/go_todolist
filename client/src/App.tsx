import Navbar from "@/components/Navbar"
import TodoForm from "@/components/TodoForm"
import TodoList from "@/components/TodoList"
import { useColorModeValue } from "@/components/ui/color-mode"
import { Container, Stack } from "@chakra-ui/react"

export const BASE_URL = "http://localhost:5000/api"
const App = () => {
  const bg = useColorModeValue("gray.", "gray.800")
  return (
    <Stack h="100vh" bg={bg}>
      <Navbar />
        <Container maxW={"600px"}>
          <TodoForm/>
          <TodoList/>
        </Container>
    </Stack>
  )
}

export default App