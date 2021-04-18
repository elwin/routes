import Master from "../../layouts/master";
import Statistics from "../../layouts/statistics";
import Notification from "../../layouts/notification";
import Table from "../../layouts/tabls";

export default function Dashboard() {
    return (
        <Master title="Dashboard">

            <Table/>

            <Notification/>

        </Master>
    )
}


