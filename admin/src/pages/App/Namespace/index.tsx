import React, { useState, useRef } from 'react';
import {
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  PlusOutlined,
  SendOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {
  CtlIngressNamespacePageData,
  CtlIngressNamespacePublish,
  CtlIngressNamespaceRemove,
  DataMapAppDataState,
} from '@/services/ant-design-pro/api';
import AddForm from './components/AddForm';
import { Button, Modal, Tag } from 'antd';
import EditForm from './components/EditForm';
import { Access, useAccess } from 'umi';
import { ShowSuccessMessage, ShowErrorMessage } from '@/services/ant-design-pro/lib';

// 数据过滤处理
const processRespData = (dataList: API.AppDataIngressNamespace[]) => {
  return dataList;
};

const DataPage: React.FC = () => {
  const [addModalVisible, setAddModalVisible] = useState<boolean>(false);
  const [editModalVisible, setEditModalVisible] = useState<boolean>(false);

  const appAccess = useAccess();

  const actionRef = useRef<ActionType>();

  const [currentRow, setCurrentRow] = useState<API.AppDataIngressNamespace>();

  const onAddModalVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setAddModalVisible(flag);
  };

  const onEditModalVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setEditModalVisible(flag);
  };

  const doPubDataAction = async (dataId?: string) => {
    const reqData: API.AppDataIdRequest = {};
    reqData.data_id = dataId;

    const respData = await CtlIngressNamespacePublish(reqData);

    if (respData.error_code === '0') {
      ShowSuccessMessage(respData.error_message);
      actionRef.current?.reload();
      return;
    }

    ShowErrorMessage(respData.error_message);
  };

  const doPubData = async (dataId?: string) => {
    Modal.confirm({
      title: '系统操作确认',
      icon: <ExclamationCircleOutlined />,
      content: '此配置数据发布操作直接影响网关结果,你确定要发布此配置吗?',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        await doPubDataAction(dataId);
      },
    });
  };

  const doRemoveDataAction = async (dataId?: string) => {
    const reqData: API.AppDataIdRequest = {};
    reqData.data_id = dataId;

    const respData = await CtlIngressNamespaceRemove(reqData);

    if (respData.error_code === '0') {
      ShowSuccessMessage(respData.error_message);
      actionRef.current?.reload();
      return;
    }

    ShowErrorMessage(respData.error_message);
  };

  const doRemoveData = async (dataId?: string) => {
    Modal.confirm({
      title: '系统操作确认',
      icon: <ExclamationCircleOutlined />,
      content: '此数据删除操作不可恢复,你确定要删除此数据吗?',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        await doRemoveDataAction(dataId);
      },
    });
  };

  // 定义列
  const columns: ProColumns<API.AppDataIngressNamespace>[] = [
    {
      title: '空间ID',
      dataIndex: 'data_id',
      hideInSearch: true,
      align: 'center',
    },
    {
      title: '空间名称',
      dataIndex: 'title',
      ellipsis: true,
    },
    {
      title: '最后更新版本',
      dataIndex: 'last_ver',
      valueType: 'dateTime',
      width: '150px',
      hideInSearch: true,
    },
    {
      title: '最后发布版本',
      dataIndex: 'last_pub_ver',
      valueType: 'dateTime',
      hideInSearch: true,
      width: '150px',
      render: (_, record) => (
        <Tag color={record.last_pub_ver == record.last_ver ? 'green' : 'red'}>
          {record.last_pub_ver}
        </Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'state',
      width: '150px',
      valueType: 'select',
      valueEnum: DataMapAppDataState,
    },
    {
      title: '操作',
      key: 'option',
      width: '200px',
      valueType: 'option',
      render: (_, record) => [
        <Access key={'access-edit'} accessible={appAccess.canAdmin}>
          <a
            key="act-edit"
            onClick={() => {
              setCurrentRow(record);
              onEditModalVisible(true);
            }}
          >
            <EditOutlined />
            &nbsp;编辑
          </a>
        </Access>,
        <a
          key="act-pub"
          onClick={() => {
            doPubData(record.data_id);
          }}
        >
          <SendOutlined />
          &nbsp;发布
        </a>,
        <Access key={'access-del'} accessible={appAccess.canAdmin}>
          <a
            key="act-del"
            onClick={() => {
              doRemoveData(record.data_id);
            }}
          >
            <DeleteOutlined />
            &nbsp;删除
          </a>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.AppDataIngressNamespace, API.AppDataIngressNamespacePageRequest>
        bordered={true}
        columns={columns}
        actionRef={actionRef}
        rowKey="data_id"
        request={CtlIngressNamespacePageData}
        postData={processRespData}
        size="small"
        toolBarRender={() => [
          <Access accessible={appAccess.canAdmin}>
            <Button
              type="primary"
              key="add-user"
              onClick={() => {
                onAddModalVisible(true);
              }}
              icon={<PlusOutlined />}
            >
              添加空间
            </Button>
          </Access>,
        ]}
      />

      <AddForm onModalVisible={onAddModalVisible} modalVisible={addModalVisible} />

      <EditForm
        onModalVisible={onEditModalVisible}
        modalVisible={editModalVisible}
        id={currentRow?.data_id}
      />
    </PageContainer>
  );
};

export default DataPage;
