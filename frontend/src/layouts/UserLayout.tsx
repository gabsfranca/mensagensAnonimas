import { MessageForm } from "../components/MessageForm";
import { StatusChecker } from "../components/StatusChecker";


export default function UserLayout() {
    return(
        <div class = "main-container">
            <MessageForm />
            <StatusChecker />
        </div>
    )
}