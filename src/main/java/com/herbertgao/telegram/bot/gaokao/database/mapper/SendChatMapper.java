package com.herbertgao.telegram.bot.gaokao.database.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.herbertgao.telegram.bot.gaokao.database.domain.SendChat;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface SendChatMapper extends BaseMapper<SendChat> {
}