import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {Pagination, Table, Tag, Typography} from 'antd';
import 'antd/dist/reset.css';
import './App.css'; // 导入自定义样式文件

const {Text, Link} = Typography;

interface DataType {
    sheet_id: string;
    title: string;
    link: string;
    release_date: string;
    abstract: string;
    keyword: string;
}

interface MetaData {
    sheet_name: string;
    update_log: string;
}

// 定义 columns 的类型
const columns = [
    {
        title: '标题',
        dataIndex: 'title',
        key: 'title',
        render: (text: string, record: DataType) => (
            <Link href={record.link} target="_blank">{text}</Link>
        ),
        width: 250, // 设置列宽
        ellipsis: true, // 超出内容使用省略号
    },
    {
        title: '链接',
        dataIndex: 'link',
        key: 'link',
        render: (text: string) => <Text copyable>{text}</Text>,
        width: 250, // 设置列宽
        ellipsis: true, // 超出内容使用省略号
    },
    {
        title: '发布日期',
        dataIndex: 'release_date',
        key: 'release_date',
        render: (date: string) => <Text>{new Date(date).toLocaleDateString()}</Text>,
        width: 150, // 设置列宽
        ellipsis: true, // 超出内容使用省略号
    },
    {
        title: '文章摘要',
        dataIndex: 'abstract',
        key: 'abstract' as string,
        render: (summary: string) => <Text ellipsis>{summary || ''}</Text>,
        width: 200, // 设置列宽，略小于标题
        ellipsis: true, // 超出内容使用省略号
    },
    {
        title: '关键字',
        dataIndex: 'keyword',
        key: 'keyword',
        render: (keywords: string) => (
            <>
                {keywords ? keywords.split(',').map((keyword: string) => (
                    <Tag color="blue" key={keyword}>{keyword}</Tag>
                )) : '-'}
            </>
        ),
        width: 200, // 设置列宽
        ellipsis: true, // 超出内容使用省略号
    },
];

const App: React.FC = () => {
    const [data, setData] = useState<DataType[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [total, setTotal] = useState<number>(0);
    const [currentPage, setCurrentPage] = useState<number>(1);
    const [sheetId, setSheetId] = useState<string>('');
    const [latestModificationTime, setLatestModificationTime] = useState<string>('');

    useEffect(() => {
        const fetchData = async () => {
            try {
                // 获取 URL 中的 sheet_id 参数
                const urlParams = new URLSearchParams(window.location.search);
                const sheetId = urlParams.get('sheet_id') || '';
                setSheetId(sheetId);

                const response = await axios.post('/api/v1/list', {
                    sheet_id: sheetId,
                    page: currentPage
                });

                // 解析数据
                if (response.data && response.data.data && response.data.data.list) {
                    setData(response.data.data.list);
                    setTotal(response.data.data.total_count); // 设置总记录数
                }
            } catch (error) {
                console.error('获取数据失败：', error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [currentPage]);

    useEffect(() => {
        const fetchMetaData = async () => {
            try {
                const response = await axios.post('/api/v1/get_meta_data', {sheet_id: sheetId});
                if (response.data) {
                    const metaData: MetaData = response.data.data;
                    setLatestModificationTime(metaData.update_log);
                }
            } catch (error) {
                console.error('获取元数据失败：', error);
            }
        };

        if (sheetId) {
            fetchMetaData();
        }
    }, [sheetId]);

    // 处理分页变化
    const handlePageChange = (page: number) => {
        setCurrentPage(page);
    };

    return (
        <div className="App">
            <div className="info-container" style={{display: 'flex', justifyContent: 'flex-end', margin: '20px'}}>
                <Text strong>飞书表格编号：</Text> <Text>{sheetId}</Text>
                <Text strong style={{marginLeft: '20px'}}>最新修改时间：</Text> <Text>{latestModificationTime}</Text>
            </div>
            <Table<DataType>
                columns={columns}
                dataSource={data.map((item, index) => ({...item, uniqueKey: `${item.sheet_id}-${index}`}))}
                rowKey="uniqueKey"
                bordered
                pagination={false}
                loading={loading}
                style={{margin: '20px', borderRadius: '8px', overflow: 'hidden'}}
                className="custom-table"
            />
            <div className="pagination-container" style={{textAlign: 'right', margin: '20px'}}>
                <Pagination
                    current={currentPage}
                    pageSize={10} // 每页显示10条数据
                    total={total}
                    onChange={handlePageChange}
                    showSizeChanger={false} // 去掉每页条数选择
                    align="center"
                />
            </div>
        </div>
    );
};

export default App;
