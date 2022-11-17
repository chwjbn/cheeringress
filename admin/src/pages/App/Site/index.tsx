import React, { useState, useRef } from 'react';
import {
  ControlOutlined,
  DeleteOutlined,
  EditOutlined,
  ExclamationCircleOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { PageContainer } from '@ant-design/pro-layout';
import type { ProColumns, ActionType } from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {
  CtlIngressSitePageData,
  CtlIngressNamespaceMap,
  DataMapAppDataIngressSiteAuthNeed,
  DataMapAppDataState,
  DataMapRuleStringOp,
  DataMapAppDataIngressActionType,
  CtlIngressSiteRemove,
} from '@/services/ant-design-pro/api';
import AddForm from './components/AddForm';
import { Button, Modal } from 'antd';
import EditForm from './components/EditForm';
import { Access, useAccess } from 'umi';
import RuleListForm from './components/RuleListForm';
import { ShowErrorMessage, ShowSuccessMessage } from '@/services/ant-design-pro/lib';

// 数据过滤处理
const processRespData = (dataList: API.AppDataIngressSite[]) => {
  return dataList;
};

const DataPage: React.FC = () => {
  const [addModalVisible, setAddModalVisible] = useState<boolean>(false);
  const [editModalVisible, setEditModalVisible] = useState<boolean>(false);
  const [xNodeListFormVisible, setNodeListFormVisible] = useState<boolean>(false);

  const appAccess = useAccess();
  const actionRef = useRef<ActionType>();

  const [currentRow, setCurrentRow] = useState<API.AppDataIngressSite>();

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

  const onNodeListFormVisible = (flag: boolean) => {
    if (!flag) {
      actionRef.current?.reload();
    }
    setNodeListFormVisible(flag);
  };

  const loadNamespaceMapData = async () => {
    const dataList: any[] = [];

    const xReq = {};

    const xRespData = await CtlIngressNamespaceMap(xReq);

    if (xRespData && xRespData.data) {
      xRespData.data.forEach((xItem) => {
        dataList.push({ label: xItem.data_name, value: xItem.data_id });
      });
    }

    return dataList;
  };

  const loadTableData = async (reqData: API.AppDataIngressSitePageRequest) => {
    const xTableData = CtlIngressSitePageData(reqData);
    return xTableData;
  };

  const doRemoveDataAction = async (dataId?: string) => {
    const reqData: API.AppDataIdRequest = {};
    reqData.data_id = dataId;

    const respData = await CtlIngressSiteRemove(reqData);

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
  const columns: ProColumns<API.AppDataIngressSite>[] = [
    {
      title: '系统ID',
      dataIndex: 'data_id',
      hideInSearch: true,
      hideInTable: true,
    },
    {
      title: '网关空间',
      dataIndex: 'namespace_id',
      valueType: 'select',
      align: 'center',
      request: loadNamespaceMapData,
      fieldProps: { showSearch: true },
    },
    {
      title: '站点名称',
      dataIndex: 'title',
      ellipsis: true,
    },
    {
      title: '域名匹配方式',
      dataIndex: 'match_op',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapRuleStringOp,
      width: '150px',
    },
    {
      title: '域名匹配内容',
      dataIndex: 'match_value',
      hideInSearch: true,
      ellipsis: true,
    },
    {
      title: '访问认证',
      dataIndex: 'auth_need',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapAppDataIngressSiteAuthNeed,
      width: '120px',
    },
    {
      title: '响应类型',
      dataIndex: 'action_type',
      valueType: 'select',
      hideInSearch: true,
      valueEnum: DataMapAppDataIngressActionType,
      width: '150px',
    },

    {
      title: '最后更新',
      dataIndex: 'update_time',
      valueType: 'dateTime',
      width: '150px',
      hideInSearch: true,
    },
    {
      title: '状态',
      dataIndex: 'state',
      valueType: 'select',
      valueEnum: DataMapAppDataState,
      width: '120px',
    },
    {
      title: '操作',
      key: 'option',
      width: '300px',
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
          key="act-rule"
          onClick={() => {
            setCurrentRow(record);
            onNodeListFormVisible(true);
          }}
        >
          <ControlOutlined />
          &nbsp;管理规则
        </a>,
        <Access key={'access-edit'} accessible={appAccess.canAdmin}>
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
      <ProTable<API.AppDataIngressSite, API.AppDataIngressSitePageRequest>
        bordered={true}
        columns={columns}
        actionRef={actionRef}
        rowKey="data_id"
        request={loadTableData}
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
              添加站点
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

      <RuleListForm
        onModalVisible={onNodeListFormVisible}
        modalVisible={xNodeListFormVisible}
        namespace_id={currentRow?.namespace_id}
        site_id={currentRow?.data_id}
      />
    </PageContainer>
  );
};

export default DataPage;
