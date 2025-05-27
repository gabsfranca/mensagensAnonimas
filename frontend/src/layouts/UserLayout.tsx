import { MessageForm } from "../components/MessageForm";
import { StatusChecker } from "../components/StatusChecker";
import { Footer } from "../components/Footer";

export default function UserLayout() {
    return(
        <div class = "page-wrapper">
            <div class = "main-container">
                <MessageForm />
                <StatusChecker />
            </div>
                <Footer />
        </div>
    )
}