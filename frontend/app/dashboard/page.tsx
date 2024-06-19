const Dashboard = () => {
    return (
        <div className="border-2 border-green-500 h-[91.5vh] flex w-full p-2">
            <div className="flex flex-col w-3/5">
                <div className="h-2/5 border-2 border-cyan-500">
                    top
                </div>
                <div className="h-3/5 border-2 border-blue-500">
                    bottom
                </div>
            </div>

            <aside className="w-2/5 border-2 border-black">
                <div className="h-full">
                    sidebar
                </div>
            </aside>
        </div>
    );
}

export default Dashboard;
