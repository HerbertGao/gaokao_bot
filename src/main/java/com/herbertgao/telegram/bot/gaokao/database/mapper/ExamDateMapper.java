package com.herbertgao.telegram.bot.gaokao.database.mapper;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.herbertgao.telegram.bot.gaokao.database.domain.ExamDate;
import org.apache.ibatis.annotations.Mapper;

@Mapper
public interface ExamDateMapper extends BaseMapper<ExamDate> {
}